package activities

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/git"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/github"
	workerframework "github.com/kube-tarian/kad/capten/common-pkg/worker-framework"
	"github.com/pkg/errors"

	"github.com/kube-tarian/kad/capten/common-pkg/plugins/argocd"
)

type ConfigureApp struct {
	config    *Config
	gitPlugin workerframework.ConfigureCICD
}

func NewConfigureApp() (*ConfigureApp, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}

	return &ConfigureApp{config: config, gitPlugin: git.New()}, nil
}

func (ca *ConfigureApp) getAccessToken(ctx context.Context, credId string) (string, error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return "", err
	}

	cred, err := credReader.GetCredential(ctx, credentials.GenericCredentialType,
		ca.config.VaultEntityName, credId)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading credential %s/%s from the vault",
			ca.config.VaultEntityName, credId)
		return "", err
	}

	return cred[accessToken], nil
}

func (ca *ConfigureApp) cloneRepos(ctx context.Context, templateRepo, customerRepo, token string) (templateDir string,
	reqRepo string, err error) {
	// Clone the template repo
	templateDir, err = os.MkdirTemp(ca.config.GitCloneDir, tmpCloneStr)
	if err != nil {
		err = fmt.Errorf("failed to create template tmp dir, err: %v", err)
		return
	}

	// Clone the customer repo
	if err = ca.gitPlugin.Clone(templateDir, templateRepo, token); err != nil {
		os.RemoveAll(templateDir)
		err = fmt.Errorf("failed to Clone template repo, err: %v", err)
		return
	}

	reqRepo, err = os.MkdirTemp(ca.config.GitCloneDir, tmpCloneStr)
	if err != nil {
		os.RemoveAll(templateDir)
		err = fmt.Errorf("failed to create tmp dir for user repo, err: %v", err)
		return
	}

	if err = ca.gitPlugin.Clone(reqRepo, customerRepo, token); err != nil {
		os.RemoveAll(templateDir)
		os.RemoveAll(reqRepo)
		err = fmt.Errorf("failed to Clone user repo, err: %v", err)
		return
	}

	return
}

func (ca *ConfigureApp) deployMainApp(ctx context.Context, fileName string) (string, string, error) {
	k8sclient, err := k8s.NewK8SClient(logging.NewLogger())
	if err != nil {
		return "", "", fmt.Errorf("failed to initalize k8s client: %v", err)
	}

	// For the testing change the reqrepo to template one
	ns, resName, err := k8sclient.DynamicClient.CreateResource(ctx, fileName)
	if err != nil {
		return "", "", fmt.Errorf("failed to create the k8s custom resource: %v", err)
	}

	return ns, resName, nil

}

func (ca *ConfigureApp) syncArgoCDApp(ctx context.Context, ns, resName string) error {
	logger.Info("RESOURCE NAME AND NAMESAPCE", resName, ns)
	client, err := argocd.NewClient(logger)
	if err != nil {
		return err
	}

	_, err = client.TriggerAppSync(ctx, ns, resName)
	if err != nil {
		return err
	}

	return nil
}

func (ca *ConfigureApp) waitForArgoCDToSync(ctx context.Context, ns, resName string) error {
	logger.Info("RESOURCE NAME AND NAMESAPCE", resName, ns)
	client, err := argocd.NewClient(logger)
	if err != nil {
		return err
	}

	for i := 0; i < 3; i++ {
		app, err := client.GetAppSyncStatus(ctx, ns, resName)
		if err != nil {
			return err
		}

		if app.Status.Sync.Status == v1alpha1.SyncStatusCodeSynced {
			break
		}

		time.Sleep(30 * time.Second)

	}

	return nil
}

func (ca *ConfigureApp) addToGit(ctx context.Context, paramType, repoUrl, token string, createPr bool) error {
	if err := ca.gitPlugin.Commit(".", "configure requested app",
		ca.config.GitDefaultCommiterName, ca.config.GitDefaultCommiterEmail); err != nil {
		return fmt.Errorf("failed to commit the changes to user repo, err: %v", err)
	}

	defaultBranch, err := ca.gitPlugin.GetDefaultBranchName()
	if err != nil {
		return fmt.Errorf("failed to get default branch of user repo, err: %v", err)
	}

	if !createPr {
		_, err = ca.createPR(ctx, repoUrl, branchName+"-"+paramType, defaultBranch, token)
		if err != nil {
			return fmt.Errorf("failed to create the PR on user repo, err: %v", err)
		}

		logger.Info("skiping push to default branch.")
		return nil
	}

	if err := ca.gitPlugin.Push(defaultBranch, token); err != nil {
		return fmt.Errorf("failed to get push to default branch, err: %v", err)
	}

	return nil
}

func (ca *ConfigureApp) createPR(ctx context.Context, repoURL, commitBranch, baseBranch, token string) (string, error) {
	op := github.NewOperation(token)
	str := strings.Split(repoURL, "/")
	return op.CreatePR(ctx, strings.TrimSuffix(str[len(str)-1], gitUrlSuffix), str[len(str)-2], "Configuring requested app", commitBranch, baseBranch, "")
}
