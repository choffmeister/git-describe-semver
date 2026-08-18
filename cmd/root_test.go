package cmd

import (
	"io/ioutil"
	"testing"

	"github.com/choffmeister/git-describe-semver/internal"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	assert := assert.New(t)
	dir, _ := ioutil.TempDir("", "example")
	author := object.Signature{Name: "Test", Email: "test@test.com"}
	_, err := run(dir, internal.GenerateVersionOptions{PrereleasePrefix: "dev"})
	assert.Error(err)

	repo, _ := git.PlainInit(dir, false)
	worktree, _ := repo.Worktree()
	_, err = run(dir, internal.GenerateVersionOptions{PrereleasePrefix: "dev"})
	assert.Error(err)

	commit1, _ := worktree.Commit("first", &git.CommitOptions{Author: &author})
	repo.CreateTag("invalid", commit1, nil)
	_, err = run(dir, internal.GenerateVersionOptions{PrereleasePrefix: "dev"})
	assert.Error(err)

	commit2, _ := worktree.Commit("first", &git.CommitOptions{Author: &author})
	repo.CreateTag("v1.0.0", commit2, nil)

	commit3, _ := worktree.Commit("second", &git.CommitOptions{Author: &author})
	result, err := run(dir, internal.GenerateVersionOptions{PrereleasePrefix: "dev"})
	assert.NoError(err)
	assert.Equal("v1.0.1-dev.1.g"+commit3.String()[0:7], *result)
}

func TestNormalizeBoolArgs(t *testing.T) {
	assert := assert.New(t)
	result := normalizeBoolArgs([]string{
		"--dir=.",
		"--fallback=0.0.0",
		"--drop-prefix=false",
		"--prerelease-timestamped=false",
		"--prerelease-prefix=pre",
		"--format=version=<version>",
	})
	assert.Equal([]string{
		"--dir=.",
		"--fallback=0.0.0",
		"--prerelease-prefix=pre",
		"--format=version=<version>",
	}, result)

	result = normalizeBoolArgs([]string{
		"--drop-prefix=true",
		"--prerelease-timestamped=1",
		"--dir=.",
	})
	assert.Equal([]string{
		"--drop-prefix",
		"--prerelease-timestamped",
		"--dir=.",
	}, result)
}
