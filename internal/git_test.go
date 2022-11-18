package internal

import (
	"io/ioutil"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/stretchr/testify/assert"
)

func TestGitTagMap(t *testing.T) {
	assert := assert.New(t)
	dir, _ := ioutil.TempDir("", "example")
	author := object.Signature{Name: "Test", Email: "test@test.com"}
	repo, _ := git.PlainInit(dir, false)
	worktree, _ := repo.Worktree()

	tags, _ := GitTagMap(*repo)
	assert.Equal(map[string]string{}, *tags)

	commit1, _ := worktree.Commit("first", &git.CommitOptions{Author: &author})
	tag1, _ := repo.CreateTag("v1.0.0", commit1, nil)
	tags, _ = GitTagMap(*repo)
	assert.Equal(commit1.String(), tag1.Hash().String())
	assert.Equal(map[string]string{
		tag1.Hash().String(): "v1.0.0",
	}, *tags)

	commit2, _ := worktree.Commit("second", &git.CommitOptions{Author: &author})
	tag2, _ := repo.CreateTag("v2.0.0", commit2, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "Foo Bar",
			Email: "foo@bar.com",
		},
		Message: "Version 2.0.0",
	})
	assert.NotEqual(commit2.String(), tag2.Hash().String())
	tags, _ = GitTagMap(*repo)
	assert.Equal(map[string]string{
		commit1.String(): "v1.0.0",
		commit2.String(): "v2.0.0",
	}, *tags)

	commit3, _ := worktree.Commit("third", &git.CommitOptions{Author: &author})
	tag3, _ := repo.CreateTag("fum", commit3, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  "Fum",
			Email: "fum@example.com",
		},
		Message: "Not a semver version tag",
	})
	assert.NotEqual(commit3.String(), tag3.Hash().String())
	tags, _ = GitTagMap(*repo)
	assert.Equal(map[string]string{
		commit1.String(): "v1.0.0",
		commit2.String(): "v2.0.0",
	}, *tags)
}

func TestGitDescribe(t *testing.T) {
	assert := assert.New(t)
	dir, _ := ioutil.TempDir("", "example")
	author := object.Signature{Name: "Test", Email: "test@test.com"}
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

	commit1, _ := worktree.Commit("first", &git.CommitOptions{Author: &author})
	test("", 1, commit1.String())

	repo.CreateTag("v1.0.0", commit1, nil)
	test("v1.0.0", 0, commit1.String())

	commit2, _ := worktree.Commit("second", &git.CommitOptions{Author: &author})
	test("v1.0.0", 1, commit2.String())

	commit3, _ := worktree.Commit("third", &git.CommitOptions{Author: &author})
	test("v1.0.0", 2, commit3.String())

	repo.CreateTag("v2.0.0", commit3, nil)
	test("v2.0.0", 0, commit3.String())
}

func TestGitDescribeWithBranch(t *testing.T) {
	assert := assert.New(t)
	dir, _ := ioutil.TempDir("", "example")
	author := object.Signature{Name: "Test", Email: "test@test.com"}
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

	commit1, _ := worktree.Commit("first", &git.CommitOptions{Author: &author})
	test("", 1, commit1.String())

	repo.CreateTag("v1.0.0", commit1, nil)
	test("v1.0.0", 0, commit1.String())

	commit2, _ := worktree.Commit("second", &git.CommitOptions{Author: &author})
	test("v1.0.0", 1, commit2.String())

	worktree.Checkout(&git.CheckoutOptions{Hash: commit1})

	commit3, _ := worktree.Commit("third", &git.CommitOptions{Author: &author})
	test("v1.0.0", 1, commit3.String())

	commit4, _ := worktree.Commit("forth", &git.CommitOptions{Author: &author, Parents: []plumbing.Hash{commit2, commit3}})
	test("v1.0.0", 2, commit4.String())
	repo.CreateTag("v2.0.0", commit3, nil)
	test("v2.0.0", 1, commit4.String())
}
