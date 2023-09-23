package activities

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kube-tarian/kad/capten/common-pkg/plugins/git"
	workerframework "github.com/kube-tarian/kad/capten/common-pkg/worker-framework"
	"github.com/kube-tarian/kad/capten/config-worker/pkg/constants"
	"github.com/kube-tarian/kad/capten/model"
)

func handleGit(ctx context.Context, params model.ConfigureParameters, payload interface{}) (model.ResponsePayload, error) {
	var err error
	respPayload := model.ResponsePayload{Status: "Failed", Message: []byte("Failed to configure the git")}
	req, ok := payload.(model.ConfigureCICD)
	if !ok {
		return respPayload, fmt.Errorf("Wrong payload: %v, recieved for configuring git", payload)
	}

	switch params.Action {
	case constants.TektonDirName:
		err = configureCICD(ctx, req, constants.TektonDirName)
	default:
		err = fmt.Errorf("unknown action %s for resouce %s", params.Action, params.Resource)
	}

	if err != nil {
		return respPayload, err
	}

	return model.ResponsePayload{Status: "Success", Message: []byte("Successfully configured the git")}, nil
}

func configureCICD(ctx context.Context, params model.ConfigureCICD, appDir string) error {
	gitPlugin := getCICDPlugin()
	configPlugin, ok := gitPlugin.(workerframework.ConfigureCICD)
	if !ok {
		return fmt.Errorf("plugin not supports Configuration for CICD activities")
	}

	dir, err := os.MkdirTemp("", "clone*")
	if err != nil {
		return err
	}

	defer os.RemoveAll(dir) // clean up

	if err := configPlugin.Clone(dir, params.RepoURL, params.Token); err != nil {
		return err
	}

	repoName := strings.Split(params.RepoURL, "/")
	// get the repoName
	cloneDir := filepath.Join(dir, strings.TrimRight(repoName[len(repoName)-1], ".git"))

	cmd := exec.Command("cp", "--recursive", filepath.Join("/", constants.GitTemplateDir, appDir), cloneDir)
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := configPlugin.Commit(appDir, fmt.Sprintf("configure %s for the repo", appDir)); err != nil {
		return err
	}

	if err := configPlugin.Push(); err != nil {
		return err
	}

	return nil
}

func getCICDPlugin() workerframework.ConfigureCICD {
	return git.New()
}
