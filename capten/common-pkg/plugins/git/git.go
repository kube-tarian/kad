package git

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/intelops/go-common/logging"
)

type GitClient struct {
	repository *git.Repository
}

func NewClient() *GitClient {
	return &GitClient{}
}

func (g *GitClient) Clone(directory, url, token string) error {
	opt := &git.CloneOptions{
		URL:             url,
		Progress:        os.Stdout,
		InsecureSkipTLS: true,
	}
	if len(token) != 0 {
		opt.Auth = &http.BasicAuth{
			Username: "dummy", // yes, this can be anything except an empty string
			Password: token,
		}
	}

	r, err := git.PlainClone(directory, false, opt)
	if err != nil {
		return err
	}

	g.repository = r
	return nil
}

func (g *GitClient) Add(path string) error {
	w, err := g.repository.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(path)
	if err != nil {
		return err
	}
	return nil
}

func (g *GitClient) Remove(path string) error {
	w, err := g.repository.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func (g *GitClient) Commit(msg, name, email string) error {
	w, err := g.repository.Worktree()
	if err != nil {
		return err
	}

	status, err := w.Status()
	if err != nil {
		return err
	}

	v, _ := json.Marshal(status)
	fmt.Println("Marshalled data => ", string(v))
	fmt.Printf("status.IsClean() => %v \n", status.IsClean())
	fmt.Printf("status.String() => %v \n", status.String())

	if status.IsClean() {
		logging.NewLogger().Info("No commit changes found.")
		return nil
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

func (g *GitClient) GetDefaultBranchName() (string, error) {
	defBranch, err := g.repository.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get the current head: %w", err)
	}

	defaultBranch := strings.Split(defBranch.String(), "/")
	return defaultBranch[len(defaultBranch)-1], nil
}

func (g *GitClient) Push(branchName, token string) error {
	defBranch, err := g.GetDefaultBranchName()
	if err != nil {
		return fmt.Errorf("failed to get the current head: %w", err)
	}

	err = g.repository.Push(&git.PushOptions{RemoteName: "origin", Force: true,
		Auth: &http.BasicAuth{
			Username: "dummy", // yes, this can be anything except an empty string
			Password: token,
		},
		InsecureSkipTLS: true,
		RefSpecs:        []config.RefSpec{config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", defBranch, branchName))}})

	return err
}
