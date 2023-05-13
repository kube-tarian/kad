package tekton

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"path/filepath"

	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func NewTekton() (*Tektonspec, error) {
	return &Tektonspec{}, nil
}

func (h *Tektonspec) Run(ctx context.Context, payload interface{}) (json.RawMessage, error) {

	request, ok := payload.(*Tektonspec)
	if !ok {
		return nil, fmt.Errorf("Error while usinn=g payload interface %v", ok)
	}

	// create new Git repository
	path := "/home/shifnazarnaz/newcap"

	filePath := "/home/shifnazarnaz/newcap/captan_new/ingress.yaml"
	//Read the above file path
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("File Reading error", err)
	}

	// Define the values for the template
	values := map[string]string{

		"Host": request.Hostname,
	}

	// Parse the template file
	tmpl, err := template.New("ingress").Parse(string(file))
	if err != nil {
		log.Fatal("Template Parsing error", err)
	}

	var buf bytes.Buffer

	// Execute the template with the provided values, writing the output to the buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		log.Fatal("Error executing template:", err)
		os.Exit(1)
	}
	// Write the modified file to the same location
	err = ioutil.WriteFile(filePath, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal("Writing file error", err)
	}
	//Initialize the repository
	repo, err := git.PlainInit(path, false)
	if err != nil {
		log.Fatalf("Error creating Git repository: %s\n", err)
	}

	// add new remote

	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",

		URLs: []string{request.Giturl},
	})
	if err != nil {
		log.Fatalf("Error creating remote: %s\n", err)

	}

	// get list of YAML files
	new := "/home/shifnazarnaz/newcap/captan_new"
	yamls, err := filepath.Glob(filepath.Join(new, "*.yaml"))

	if err != nil {
		log.Fatalf("Error getting YAML files: %s\n", err)

	}
	fmt.Println("Yaml File is", yamls)

	//authentication
	auth := &http.BasicAuth{
		Username: "Shifna12Zarnaz",
		Password: "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",
	}

	// add all YAML files to Git index
	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Error getting worktree: %s\n", err)

	}
	_, err = worktree.Add(".")
	if err != nil {
		log.Fatalf("Error adding files to index: %s\n", err)

	}

	// commit changes
	commit, err := worktree.Commit("Add YAML files", &git.CommitOptions{
		Author: &object.Signature{
			Name:  request.Name,
			Email: request.Mail,
		},
	})
	if err != nil {
		log.Fatalf("Error committing changes: %s\n", err)

	}
	fmt.Println("Commit is", commit)
	// push changes to remote repository
	err = remote.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{"refs/heads/master:refs/heads/master"},
		Auth:       auth,
	})
	if err != nil {
		log.Fatalf("Error while pushing the file: %s\n", err)

	}
	return nil, nil
}
