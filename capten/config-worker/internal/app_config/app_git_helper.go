package appconfig

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	"github.com/kube-tarian/kad/capten/common-pkg/plugins/git"
	"github.com/pkg/errors"

	"github.com/kube-tarian/kad/capten/common-pkg/plugins/argocd"
)

const (
	tmpGitProjectCloneStr          = "clone*"
	gitProjectAccessTokenAttribute = "accessToken"
	gitProjectUserId               = "userID"
	gitUrlSuffix                   = ".git"
	kubeConfig                     = "kubeconfig"
	k8sEndpoint                    = "endpoint"
	k8sClusterCA                   = "clusterCA"
	cosignKey                      = "cosign.key"
	cosignPub                      = "cosign.pub"
)

type Config struct {
	GitDefaultCommitMessage  string `envconfig:"GIT_COMMIT_MSG" default:"capten-config-update"`
	GitDefaultCommiterName   string `envconfig:"GIT_COMMIT_NAME" default:"capten-bot"`
	GitDefaultCommiterEmail  string `envconfig:"GIT_COMMIT_EMAIL" default:"capten-bot@intelops.dev"`
	GitVaultEntityName       string `envconfig:"GIT_VAULT_ENTITY_NAME" default:"git-project"`
	GitCloneDir              string `envconfig:"GIT_CLONE_DIR" default:"/gitCloneDir"`
	GitBranchName            string `envconfig:"GIT_BRANCH_NAME" default:"capten-template-bot"`
	ManagedClusterEntityName string `envconfig:"MANAGED_CLUSER_VAULT_ENTITY_NAME" default:"managedcluster"`
}

var logger = logging.NewLogger()

type AppGitConfigHelper struct {
	cfg         Config
	gitClient   *git.GitClient
	accessToken string
}

func NewAppGitConfigHelper() (*AppGitConfigHelper, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &AppGitConfigHelper{cfg: cfg, gitClient: git.NewClient()}, nil
}

func (ca *AppGitConfigHelper) GetGitCreds(ctx context.Context, projectId string) (string, string, error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return "", "", err
	}

	cred, err := credReader.GetCredential(ctx, credentials.GenericCredentialType,
		ca.cfg.GitVaultEntityName, projectId)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading credential %s/%s from the vault",
			ca.cfg.GitVaultEntityName, projectId)
		return "", "", err
	}

	return cred[gitProjectUserId], cred[gitProjectAccessTokenAttribute], nil
}

func (ca *AppGitConfigHelper) GetClusterCreds(ctx context.Context, entityName, projectId string) (string, string, string, error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return "", "", "", err
	}

	cred, err := credReader.GetCredential(ctx, credentials.GenericCredentialType,
		entityName, projectId)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading credential %s/%s from the vault",
			entityName, projectId)
		return "", "", "", err
	}

	return cred[kubeConfig], cred[k8sClusterCA], cred[k8sEndpoint], nil
}

func (ca *AppGitConfigHelper) GetCosingKeys(ctx context.Context, entityName, projectId string) (string, string, error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return "", "", err
	}

	cred, err := credReader.GetCredential(ctx, credentials.GenericCredentialType,
		entityName, projectId)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading credential %s/%s from the vault",
			entityName, projectId)
		return "", "", err
	}

	return cred[cosignKey], cred[cosignPub], nil
}

func (ca *AppGitConfigHelper) GetContainerRegCreds(ctx context.Context, entityName, projectId string) (string, string, error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return "", "", err
	}

	cred, err := credReader.GetCredential(ctx, credentials.GenericCredentialType,
		entityName, projectId)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading credential %s/%s from the vault",
			entityName, projectId)
		return "", "", err
	}

	return cred["username"], cred["password"], nil
}

func (ca *AppGitConfigHelper) CloneTemplateRepo(ctx context.Context, repoURL, projectId string) (templateDir string, err error) {
	templateDir, err = os.MkdirTemp(ca.cfg.GitCloneDir, tmpGitProjectCloneStr)
	if err != nil {
		err = fmt.Errorf("failed to create template tmp dir, err: %v", err)
		return
	}

	gitClient := git.NewClient()
	if err = gitClient.Clone(templateDir, repoURL, ""); err != nil {
		os.RemoveAll(templateDir)
		err = fmt.Errorf("failed to Clone template repo, err: %v", err)
		return
	}
	return
}

func (ca *AppGitConfigHelper) CloneUserRepo(ctx context.Context, repoURL, projectId string) (reqRepo string, err error) {
	_, accessToken, err := ca.GetGitCreds(ctx, projectId)
	if err != nil {
		err = fmt.Errorf("failed to get token from vault, %v", err)
		return
	}
	ca.accessToken = accessToken

	reqRepo, err = os.MkdirTemp(ca.cfg.GitCloneDir, tmpGitProjectCloneStr)
	if err != nil {
		err = fmt.Errorf("failed to create tmp dir for user repo, err: %v", err)
		return
	}

	if err = ca.gitClient.Clone(reqRepo, repoURL, accessToken); err != nil {
		os.RemoveAll(reqRepo)
		err = fmt.Errorf("failed to Clone user repo, %v", err)
		return
	}
	return
}

func (ca *AppGitConfigHelper) DeployMainApp(ctx context.Context, fileName string) (string, string, error) {
	k8sclient, err := k8s.NewK8SClient(logging.NewLogger())
	if err != nil {
		return "", "", fmt.Errorf("failed to initalize k8s client, %v", err)
	}

	ns, resName, err := k8sclient.DynamicClient.CreateResourceFromFile(ctx, fileName)
	if err != nil {
		return "", "", fmt.Errorf("failed to create the k8s custom resource: %v", err)
	}

	return ns, resName, nil

}

func (ca *AppGitConfigHelper) SyncArgoCDApp(ctx context.Context, ns, resName string) error {
	client, err := argocd.NewClient(logger)
	if err != nil {
		fmt.Println("Client error ")
		fmt.Println(err.Error())
		return err
	}

	_, err = client.TriggerAppSync(ctx, ns, resName)
	if err != nil {
		fmt.Println("TriggerAppSync")
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (ca *AppGitConfigHelper) DeleteArgoCDApp(ctx context.Context, ns, resName, mainApp string) error {
	client, err := argocd.NewClient(logger)
	if err != nil {
		return err
	}

	// _, err = client.Delete(&model.DeleteRequestPayload{Namespace: ns, ReleaseName: resName})
	// if err != nil {
	// 	return err
	// }

	_, err = client.TriggerAppSync(ctx, ns, mainApp)
	if err != nil {
		return err
	}

	return nil
}

func (ca *AppGitConfigHelper) CreateCluster(ctx context.Context, id, clusterName string) (string, error) {
	credReader, err := credentials.NewCredentialReader(ctx)
	if err != nil {
		err = errors.WithMessage(err, "error in initializing credential reader")
		return "", err
	}

	cred, err := credReader.GetCredential(ctx, credentials.GenericCredentialType,
		ca.cfg.ManagedClusterEntityName, id)
	if err != nil {
		err = errors.WithMessagef(err, "error while reading credential %s/%s from the vault",
			ca.cfg.GitVaultEntityName, id)
		return "", err
	}

	client, err := argocd.NewClient(logger)
	if err != nil {
		return "", err
	}

	err = client.CreateOrUpdateCluster(ctx, clusterName, cred[kubeConfig])
	if err != nil {
		return "", err
	}

	return cred[k8sEndpoint], nil
}

func (ca *AppGitConfigHelper) WaitForArgoCDToSync(ctx context.Context, ns, resName string) error {
	client, err := argocd.NewClient(logger)
	if err != nil {
		return err
	}

	synched := false
	for i := 0; i < 3; i++ {
		app, err := client.GetAppSyncStatus(ctx, ns, resName)
		if err != nil {
			return fmt.Errorf("app %s synch staus fetch failed", resName)
		}

		if app.Status.Sync.Status == v1alpha1.SyncStatusCodeSynced {
			synched = true
			break
		}

		time.Sleep(30 * time.Second)
	}

	if !synched {
		return fmt.Errorf("app %s not synched", resName)
	}
	return nil
}

func (ca *AppGitConfigHelper) AddFilesToRepo(paths []string) error {
	for _, path := range paths {
		if err := ca.gitClient.Add(path); err != nil {
			return fmt.Errorf("failed to add '%s' the changes to repo, %v", path, err)
		}
	}
	return nil
}

func (ca *AppGitConfigHelper) RemoveFilesFromRepo(paths []string) error {
	for _, path := range paths {
		if err := ca.gitClient.Remove(path); err != nil {
			if !strings.Contains(err.Error(), "entry not found") {
				return fmt.Errorf("failed to remove '%s' the changes from repo, %v", path, err)
			}
		}
	}
	return nil
}

func (ca *AppGitConfigHelper) CommitRepoChanges() error {
	if len(ca.accessToken) == 0 {
		return fmt.Errorf("git project access token empty")
	}

	if ok, err := ca.gitClient.IsEmptyChanges(); ok {
		logging.NewLogger().Info("Empty commit changes, %v", err)
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check for empty changes, %v", err)
	}

	if err := ca.gitClient.Commit(ca.cfg.GitDefaultCommitMessage,
		ca.cfg.GitDefaultCommiterName, ca.cfg.GitDefaultCommiterEmail); err != nil {
		return fmt.Errorf("failed to commit the changes to user repo, %v", err)
	}

	defaultBranch, err := ca.gitClient.GetDefaultBranchName()
	if err != nil {
		return fmt.Errorf("failed to get default branch of user repo, %v", err)
	}

	if err := ca.gitClient.Push(defaultBranch, ca.accessToken); err != nil {
		return fmt.Errorf("failed to get push to default branch, %v", err)
	}
	return nil
}
