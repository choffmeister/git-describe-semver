package main

import (
	"io/ioutil"
	"testing"

	"github.com/choffmeister/git-describe-semver/internal"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	assert := assert.New(t)
	dir, _ := ioutil.TempDir("", "example")
	_, err := Run(dir, internal.GenerateVersionOptions{})
	assert.Error(err)

	repo, _ := git.PlainInit(dir, false)
	worktree, _ := repo.Worktree()
	_, err = Run(dir, internal.GenerateVersionOptions{})
	assert.Error(err)

	commit1, _ := worktree.Commit("first", &git.CommitOptions{})
	repo.CreateTag("invalid", commit1, nil)
	_, err = Run(dir, internal.GenerateVersionOptions{})
	assert.Error(err)

	commit2, _ := worktree.Commit("first", &git.CommitOptions{})
	repo.CreateTag("v1.0.0", commit2, nil)

	commit3, _ := worktree.Commit("second", &git.CommitOptions{})
	result, err := Run(dir, internal.GenerateVersionOptions{})
	assert.NoError(err)
	assert.Equal("v1.0.1-dev.1.g"+commit3.String()[0:7], *result)
}
