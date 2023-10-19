package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/intelops/go-common/credentials"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/git"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/github"
	workerframework "github.com/kube-tarian/kad/capten/common-pkg/worker-framework"
	"github.com/kube-tarian/kad/capten/model"
	cp "github.com/otiai10/copy"
	"github.com/pkg/errors"
)

func handleGit(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (model.ResponsePayload, error) {
	var err error
	respPayload := model.ResponsePayload{Status: "Failed", Message: json.RawMessage("{\"error\": \"requested payload is wrong\"}")}
	req := &model.UseCase{}
	err = json.Unmarshal(payload, req)
	if err != nil {
		return respPayload, fmt.Errorf("Wrong payload: %v, recieved for configuring git", payload)
	}

	// read from the vault
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return model.ResponsePayload{Status: "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err))}, err
	}

	cred, err := credReader.GetCredential(ctx, credentials.GenericCredentialType,
		req.VaultCredIdentifier, req.VaultCredIdentifier)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading credential %s/%s from the vault",
			req.VaultCredIdentifier, req.VaultCredIdentifier)
		return model.ResponsePayload{Status: "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err))}, err
	}

	switch req.Type {
	case "tekton":
		err = configureCICD(ctx, req, cred["accessToken"])
		// Once we finalize what needs to be replaced then we can come and work here.
	default:
		err = fmt.Errorf("unknown use case type %s for resouce", req.Type)
	}

	if err != nil {
		fmt.Println("ERROR: ", err)
		return model.ResponsePayload{Status: "Failed"}, err
	}

	return model.ResponsePayload{Status: "Success"}, nil
}

func configureCICD(ctx context.Context, params *model.UseCase, token string) error {
	config, _ := GetConfig()
	gitPlugin := getCICDPlugin()
	configPlugin, ok := gitPlugin.(workerframework.ConfigureCICD)
	if !ok {
		return fmt.Errorf("plugin not supports Configuration for CICD activities")
	}

	// Clone the template repo
	templateDir, err := os.MkdirTemp(config.GitCLoneDir, "clone*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(templateDir)

	if err := configPlugin.Clone(templateDir, config.CiCDTemplateRepo, token); err != nil {
		return err
	}

	reqRepo, err := os.MkdirTemp(config.GitCLoneDir, "clone*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(reqRepo) // clean up

	if err := configPlugin.Clone(reqRepo, params.RepoURL, token); err != nil {
		return err
	}

	for _, dir := range strings.Split(params.Type, ",") {
		err = cp.Copy(filepath.Join(templateDir, dir), filepath.Join(reqRepo, dir))
		if err != nil {
			return err
		}

	}

	if err := configPlugin.Commit(".", "configure CICD for the repo",
		config.GitDefaultCommiterName, config.GitDefaultCommiterEmail); err != nil {
		return err
	}

	if err := configPlugin.Push(branchName+"-"+params.Type, token); err != nil {
		return err
	}

	defaultBranch, err := configPlugin.GetDefaultBranchName()
	if err != nil {
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
