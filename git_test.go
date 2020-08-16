package main

import (
	"io/ioutil"
	"testing"

	"github.com/go-git/go-git/v5"

	"github.com/stretchr/testify/assert"
)

func TestGitTagMap(t *testing.T) {
	assert := assert.New(t)
	dir, _ := ioutil.TempDir("", "example")
	repo, _ := git.PlainInit(dir, false)
	worktree, _ := repo.Worktree()

	tags, _ := GitTagMap(*repo)
	assert.Equal(map[string]string{}, *tags)

	commit1, _ := worktree.Commit("first", &git.CommitOptions{})
	tag1, _ := repo.CreateTag("v1.0.0", commit1, nil)
	tags, _ = GitTagMap(*repo)
	assert.Equal(map[string]string{
		tag1.Hash().String(): "v1.0.0",
	}, *tags)

	commit2, _ := worktree.Commit("second", &git.CommitOptions{})
	tag2, _ := repo.CreateTag("v2.0.0", commit2, nil)
	tags, _ = GitTagMap(*repo)
	assert.Equal(map[string]string{
		tag1.Hash().String(): "v1.0.0",
		tag2.Hash().String(): "v2.0.0",
	}, *tags)
}

func TestGitDescribe(t *testing.T) {
	assert := assert.New(t)
	dir, _ := ioutil.TempDir("", "example")
	repo, _ := git.PlainInit(dir, false)
	worktree, _ := repo.Worktree()
	_, _, _, err := GitDescribe(*repo)
	assert.Error(err)
	test := func(expectedTagName string, expectedCounter int, expectedHeadHash string) {
		actualTagName, actualCounter, actualHeadHash, err := GitDescribe(*repo)
		assert.NoError(err)
		assert.Equal(expectedTagName, *actualTagName)
		assert.Equal(expectedCounter, *actualCounter)
		assert.Equal(expectedHeadHash, *actualHeadHash)
	}

	commit1, _ := worktree.Commit("first", &git.CommitOptions{})
	test("", 1, commit1.String())

	repo.CreateTag("v1.0.0", commit1, nil)
	test("v1.0.0", 0, commit1.String())

	commit2, _ := worktree.Commit("second", &git.CommitOptions{})
	test("v1.0.0", 1, commit2.String())

	commit3, _ := worktree.Commit("second", &git.CommitOptions{})
	test("v1.0.0", 2, commit3.String())

	repo.CreateTag("v2.0.0", commit3, nil)
	test("v2.0.0", 0, commit3.String())
}
