package main

import (
	"io/ioutil"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateVersion(t *testing.T) {
	assert := assert.New(t)
	test := func(inputTagName string, inputCounter int, inputHeadHash string, inputOpts GenerateVersionOptions, expected string) {
		actual, err := GenerateVersion(inputTagName, inputCounter, inputHeadHash, inputOpts)
		if assert.NoError(err) {
			assert.Equal(expected, *actual)

		}
	}

	test("0.0.0", 0, "abc1234", GenerateVersionOptions{}, "0.0.0")
	test("0.0.0", 1, "abc1234", GenerateVersionOptions{}, "0.0.1-dev.1.gabc1234")
	test("0.0.0-rc1", 1, "abc1234", GenerateVersionOptions{}, "0.0.0-rc1.dev.1.gabc1234")
	test("0.0.0-rc.1", 1, "abc1234", GenerateVersionOptions{}, "0.0.0-rc.1.dev.1.gabc1234")
	test("0.0.0-rc.1+foobar", 1, "abc1234", GenerateVersionOptions{}, "0.0.0-rc.1.dev.1.gabc1234+foobar")
	test("v0.0.0-rc.1+foobar", 1, "abc1234", GenerateVersionOptions{}, "v0.0.0-rc.1.dev.1.gabc1234+foobar")
	test("v0.0.0-rc.1+foobar", 1, "abc1234", GenerateVersionOptions{DropTagNamePrefix: true}, "0.0.0-rc.1.dev.1.gabc1234+foobar")

	test("", 1, "abc1234", GenerateVersionOptions{FallbackTagName: "0.0.0"}, "0.0.0-dev.1.gabc1234")
	test("", 1, "abc1234", GenerateVersionOptions{FallbackTagName: "v0.0.0"}, "v0.0.0-dev.1.gabc1234")
	test("", 1, "abc1234", GenerateVersionOptions{FallbackTagName: "v0.0.0", DropTagNamePrefix: true}, "0.0.0-dev.1.gabc1234")
	_, err := GenerateVersion("", 1, "abc1234", GenerateVersionOptions{})
	assert.Error(err)
}

func TestRun(t *testing.T) {
	assert := assert.New(t)
	dir, _ := ioutil.TempDir("", "example")
	_, err := Run(dir, GenerateVersionOptions{})
	assert.Error(err)

	repo, _ := git.PlainInit(dir, false)
	worktree, _ := repo.Worktree()
	_, err = Run(dir, GenerateVersionOptions{})
	assert.Error(err)

	commit1, _ := worktree.Commit("first", &git.CommitOptions{})
	repo.CreateTag("invalid", commit1, nil)
	_, err = Run(dir, GenerateVersionOptions{})
	assert.Error(err)

	commit2, _ := worktree.Commit("first", &git.CommitOptions{})
	repo.CreateTag("v1.0.0", commit2, nil)

	commit3, _ := worktree.Commit("second", &git.CommitOptions{})
	result, err := Run(dir, GenerateVersionOptions{})
	assert.NoError(err)
	assert.Equal("v1.0.1-dev.1.g"+commit3.String()[0:7], *result)
}

func TestMain(t *testing.T) {
	main()
}
