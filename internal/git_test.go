package internal

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func setUpDotGitDirTest(assert *assert.Assertions) (string, string) {
	testDir, err := os.MkdirTemp("", "test")
	assert.NoError(err, "failed to create temp dir")

	gitDirPath := filepath.Join(testDir, GitDirName)
	err = os.Mkdir(gitDirPath, 0750)
	assert.NoError(err, "failed to create git dir")

	return testDir, gitDirPath
}

func setUpDotGitFileTest(assert *assert.Assertions) (string, string, string) {
	testDir, err := os.MkdirTemp("", "test")
	assert.NoError(err, "failed to create temp dir")

	actualDotGitPath := filepath.Join(testDir, "actual")
	err = os.Mkdir(actualDotGitPath, 0750)
	assert.NoError(err, "failed to create actual git dir")

	wtPath := filepath.Join(testDir, "my_worktree")
	err = os.Mkdir(wtPath, 0750)
	assert.NoError(err, "failed to create worktree dir")

	wtDotGitPath := filepath.Join(wtPath, GitDirName)
	contents := GitDirPrefix + actualDotGitPath
	err = os.WriteFile(wtDotGitPath, []byte(contents), 0666)
	assert.NoError(err, "failed to write git dir file in worktree")

	return testDir, actualDotGitPath, wtPath
}

func TestFindGitDir(t *testing.T) {
	t.Run(".git is a directory", func(t *testing.T) {
		assert := assert.New(t)
		testDir, gitDirPath := setUpDotGitDirTest(assert)
		defer os.RemoveAll(testDir)

		result, err := FindGitDir(testDir)
		assert.NoError(err, "failed to find git dir")
		assert.Equal(gitDirPath, result)
	})
	t.Run(".git is a file pointing to another directory", func(t *testing.T) {
		assert := assert.New(t)
		testDir, actualDotGitPath, wtPath := setUpDotGitFileTest(assert)
		defer os.RemoveAll(testDir)

		result, err := FindGitDir(wtPath)
		assert.NoError(err, "failed to find git dir in worktree")
		assert.Equal(actualDotGitPath, result)
	})
}

func TestShouldEnableCommonDir(t *testing.T) {
	t.Run(".git is a directory", func(t *testing.T) {
		assert := assert.New(t)
		testDir, gitDirPath := setUpDotGitDirTest(assert)
		defer os.RemoveAll(testDir)

		result, err := shouldEnableCommondDir(gitDirPath)
		assert.NoError(err, "failed evaluating whether to enable commond dir")
		assert.False(result)
	})
	t.Run(".git is a file pointing to another directory", func(t *testing.T) {
		assert := assert.New(t)
		testDir, actualDotGitPath, _ := setUpDotGitFileTest(assert)
		defer os.RemoveAll(testDir)

		cdPath := filepath.Join(actualDotGitPath, CommonDirName)
		contents := "../my_worktree"
		err := os.WriteFile(cdPath, []byte(contents), 0666)
		assert.NoError(err, "failed writing commondir file")

		result, err := shouldEnableCommondDir(actualDotGitPath)
		assert.NoError(err, "failed evaluating whether to enable commond dir")
		assert.True(result)
	})
}
