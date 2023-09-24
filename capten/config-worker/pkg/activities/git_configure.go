package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kube-tarian/kad/capten/common-pkg/plugins/git"
	workerframework "github.com/kube-tarian/kad/capten/common-pkg/worker-framework"
	"github.com/kube-tarian/kad/capten/model"
)

func handleGit(ctx context.Context, params model.ConfigureParameters, payload json.RawMessage) (model.ResponsePayload, error) {
	var err error
	respPayload := model.ResponsePayload{Status: "Failed", Message: json.RawMessage("{\"error\": \"requested payload is wrong\"}")}
	req := &model.ConfigureCICD{}
	err = json.Unmarshal(payload, req)
	if err != nil {
		return respPayload, fmt.Errorf("Wrong payload: %v, recieved for configuring git", payload)
	}

	switch params.Action {
	case TektonDirName:
		err = configureCICD(ctx, req, TektonDirName)
		// Raise a PR
	default:
		err = fmt.Errorf("unknown action %s for resouce %s", params.Action, params.Resource)
	}

	if err != nil {
		return model.ResponsePayload{Status: "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", err))}, err
	}

	return model.ResponsePayload{Status: "Success"}, nil
}

func configureCICD(ctx context.Context, params *model.ConfigureCICD, appDir string) error {
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

	cmd := exec.Command("cp", "--recursive", filepath.Join("./", GitTemplateDir, appDir), dir)
	if err := cmd.Run(); err != nil {
		return err
	}

	if err := configPlugin.Commit(appDir, fmt.Sprintf("configure %s for the repo", appDir)); err != nil {
		return err
	}

	if err := configPlugin.Push(appDir+"-"+branchSuffix, params.Token); err != nil {
		return err
	}

	return nil
}

func getCICDPlugin() workerframework.ConfigureCICD {
	return git.New()
}
