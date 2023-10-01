package git

import (
	"fmt"
	"os"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Operation struct {
	repository *git.Repository
}

func New() *Operation {
	return &Operation{}
}

func (op *Operation) Clone(directory, url, token string) error {
	fmt.Println("CLONE TOEKN: ", token)
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "dummy", // yes, this can be anything except an empty string
			Password: token,
		},
		URL:             url,
		Progress:        os.Stdout,
		InsecureSkipTLS: true,
	})

	if err != nil {
		return err
	}

	op.repository = r

	return nil
}

func (op *Operation) Commit(path, msg, name, email string) error {
	w, err := op.repository.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(path)
	if err != nil {
		return err
	}

	_, err = w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  name,
			Email: email,
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (op *Operation) GetDefaultBranchName() (string, error) {
	defBranch, err := op.repository.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get the current head: %w", err)
	}

	return string(defBranch.Name()), nil
}

func (op *Operation) Push(branchName, token string) error {
	defBranch, err := op.GetDefaultBranchName()
	if err != nil {
		return fmt.Errorf("failed to get the current head: %w", err)
	}

	return op.repository.Push(&git.PushOptions{RemoteName: "origin", Force: true,
		Auth: &http.BasicAuth{
			Username: "dummy", // yes, this can be anything except an empty string
			Password: token,
		},
		InsecureSkipTLS: true,
		RefSpecs:        []config.RefSpec{config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", defBranch, branchName))}})
}
