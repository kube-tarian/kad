package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/intelops/go-common/credentials"
	agentmodel "github.com/kube-tarian/kad/capten/agent/pkg/model"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/git"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/github"
	workerframework "github.com/kube-tarian/kad/capten/common-pkg/worker-framework"
	"github.com/kube-tarian/kad/capten/model"
	cp "github.com/otiai10/copy"
	"github.com/pkg/errors"
)

type HandleGit struct {
	config       *Config
	pluginConfig *PluginConfigExtractor
}

func NewHandleGit() (*HandleGit, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	pluginConfig, err := NewPluginExtractor(
		config.TektonPluginConfig,
		config.CrossPlanePluginConfig,
		config.CrossPlaneProviderPluginConfig,
	)
	if err != nil {
		return nil, err
	}

	return &HandleGit{config: config, pluginConfig: pluginConfig}, nil
}

func (hg *HandleGit) handleGit(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (model.ResponsePayload, error) {
	var err error
	respPayload := model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed), Message: json.RawMessage("{\"error\": \"requested payload is wrong\"}")}
	req := &model.UseCase{}
	err = json.Unmarshal(payload, req)
	if err != nil {
		return respPayload, fmt.Errorf("Wrong payload: %v, recieved for configuring git", payload)
	}

	// read from the vault
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed),
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err))}, err
	}

	cred, err := credReader.GetCredential(ctx, credentials.GenericCredentialType,
		hg.config.VaultEntityName, req.VaultCredIdentifier)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading credential %s/%s from the vault",
			hg.config.VaultEntityName, req.VaultCredIdentifier)
		return model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed),
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err))}, err
	}

	switch req.Type {
	case Tekton:
		err = hg.configureCICD(ctx, req, hg.pluginConfig.tektonGetGitRepo(),
			hg.pluginConfig.tektonGetGitConfigPath(), cred["accessToken"])
	case CrossPlane:
		if err = hg.configureCICD(ctx, req, hg.pluginConfig.crossplaneGetGitRepo(),
			hg.pluginConfig.crossplaneGetGitConfigPath(), cred["accessToken"]); err != nil {
			fmt.Println("ERROR while configureCICD: ", err)
			return model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed)}, err
		}
		err = hg.configureCrossplaneProvider(ctx, req,
			hg.pluginConfig.crossplaneGetConfigMainApp(), cred["accessToken"])
	}
	// Once we finalize what needs to be replaced then we can come and work here.

	if err != nil {
		fmt.Println("ERROR: ", err)
		return model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusFailed)}, err
	}

	return model.ResponsePayload{Status: string(agentmodel.WorkFlowStatusCompleted)}, nil
}

func (hg *HandleGit) configureCICD(ctx context.Context, params *model.UseCase, templateRepo, pathInRepo, token string) error {
	gitPlugin := getCICDPlugin()
	configPlugin, ok := gitPlugin.(workerframework.ConfigureCICD)
	if !ok {
		return fmt.Errorf("plugin not supports Configuration for CICD activities")
	}

	// Clone the template repo
	templateDir, err := os.MkdirTemp(hg.config.GitCLoneDir, "clone*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(templateDir)

	if err := configPlugin.Clone(templateDir, templateRepo, token); err != nil {
		return err
	}

	reqRepo, err := os.MkdirTemp(hg.config.GitCLoneDir, "clone*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(reqRepo) // clean up

	if err := configPlugin.Clone(reqRepo, params.RepoURL, token); err != nil {
		return err
	}

	for _, dir := range strings.Split(pathInRepo, ",") {
		err = cp.Copy(filepath.Join(templateDir, dir), filepath.Join(reqRepo, dir))
		if err != nil {
			return err
		}

	}

	if err := configPlugin.Commit(".", "configure CICD for the repo",
		hg.config.GitDefaultCommiterName, hg.config.GitDefaultCommiterEmail); err != nil {
		return err
	}

	localBranchName := branchName + "-" + params.Type
	defaultBranch, err := configPlugin.GetDefaultBranchName()
	if err != nil {
		return err
	}

	if params.PushToDefaultBranch {
		localBranchName = defaultBranch
	}

	if err := configPlugin.Push(localBranchName, token); err != nil || params.PushToDefaultBranch {
		return err
	}

	_, err = createPR(ctx, params.RepoURL, branchName+"-"+params.Type, defaultBranch, token)
	if err != nil {
		return err
	}

	return nil
}

func getCICDPlugin() workerframework.ConfigureCICD {
	return git.New()
}

func createPR(ctx context.Context, repoURL, commitBranch, baseBranch, token string) (string, error) {
	op := github.NewOperation(token)
	str := strings.Split(repoURL, "/")
	return op.CreatePR(ctx, strings.TrimSuffix(str[len(str)-1], gitUrlSuffix), str[len(str)-2], "Configuring CI/CD", commitBranch, baseBranch, "")
}
