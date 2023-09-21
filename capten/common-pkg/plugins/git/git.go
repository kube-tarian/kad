package git

import (
	"os"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type Operation struct {
}

func (op *Operation) Clone(directory, url, token string) (*git.Repository, error) {
	r, err := git.PlainClone(directory, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: "dummy", // yes, this can be anything except an empty string
			Password: token,
		},
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		return nil, err
	}

	return r, err
}

func (op *Operation) Commit(r *git.Repository, path, msg string) error {
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(path)
	if err != nil {
		return err
	}

	_, err = w.Commit("example go-git commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "capten-agent-bot",
			Email: "capten-agent-bot@intelops.dev",
			When:  time.Now(),
		},
	})

	if err != nil {
		return err
	}

	return nil
}

func (op *Operation) Push(r *git.Repository) error {
	return r.Push(&git.PushOptions{})
}
