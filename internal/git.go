package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	// Prefix found in .git files that point to another location
	GitDirPrefix = "gitdir: "

	GitDirName    = ".git"
	CommonDirName = "commondir"
)

// GitTagMap ...
func GitTagMap(repo git.Repository) (*map[string]string, error) {
	iter, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	tagMap := map[string]string{}
	err = iter.ForEach(func(r *plumbing.Reference) error {
		tag, _ := repo.TagObject(r.Hash())
		if SemVerParse(r.Name().Short()) == nil {
			// Filter out tags that are not semver
			return nil
		}
		if tag == nil {
			tagMap[r.Hash().String()] = r.Name().Short()
		} else {
			c, err := tag.Commit()
			if err != nil {
				return err
			}
			tagMap[c.Hash.String()] = r.Name().Short()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &tagMap, nil
}

// GitDescribe ...
func GitDescribe(repo git.Repository) (*string, *int, *string, error) {
	type gitDescribeNode struct {
		Commit   object.Commit
		Distance int
	}

	head, err := repo.Head()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to find head: %v", err)
	}
	headHash := head.Hash().String()
	tags, err := GitTagMap(repo)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to get tags: %v", err)
	}
	commits, err := repo.Log(&git.LogOptions{
		From:  head.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to get log: %v", err)
	}
	state := map[string]gitDescribeNode{}
	counter := 0
	tagHash := ""
	commits.ForEach(func(c *object.Commit) error {
		node, found := state[c.Hash.String()]
		if !found {
			node = gitDescribeNode{
				Commit:   *c,
				Distance: 0,
			}
			state[c.Hash.String()] = node
		}
		c.Parents().ForEach(func(p *object.Commit) error {
			_, found := state[p.Hash.String()]
			if !found {
				state[p.Hash.String()] = gitDescribeNode{
					Commit:   *p,
					Distance: node.Distance + 1,
				}
			}
			return nil
		})

		_, foundTag := (*tags)[c.Hash.String()]
		if tagHash == "" && foundTag {
			counter = state[c.Hash.String()].Distance
			tagHash = c.Hash.String()
		}
		return nil
	})
	if tagHash == "" {
		for _, node := range state {
			if node.Distance+1 > counter {
				counter = node.Distance + 1
			}
		}
		tagName := ""
		return &tagName, &counter, &headHash, nil
	}
	tagName := (*tags)[tagHash]
	return &tagName, &counter, &headHash, nil
}

func OpenRepository(dir string) (*git.Repository, error) {
	gitDir, err := FindGitDir(dir)
	if err != nil {
		return nil, err
	}
	enableCommonDir, err := shouldEnableCommondDir(gitDir)
	if err != nil {
		return nil, err
	}
	openOpts := &git.PlainOpenOptions{EnableDotGitCommonDir: enableCommonDir}
	return git.PlainOpenWithOptions(dir, openOpts)
}

func shouldEnableCommondDir(gitDir string) (bool, error) {
	cdPath := filepath.Join(gitDir, CommonDirName)
	st, err := os.Stat(cdPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	if st.IsDir() {
		return false, fmt.Errorf("expected to be a file, not directory: %s", cdPath)
	}
	return true, nil
}

func FindGitDir(dir string) (string, error) {
	gitDirPath := filepath.Join(dir, GitDirName)
	st, err := os.Stat(gitDirPath)
	if err != nil {
		return "", err
	}
	if st.IsDir() {
		return gitDirPath, nil
	}
	// It is a file, read the contents
	contents, err := os.ReadFile(gitDirPath)
	if err != nil {
		return "", err
	}

	line := string(contents)
	if !strings.HasPrefix(line, GitDirPrefix) {
		return "", fmt.Errorf(".git file has no %s prefix", GitDirPrefix)
	}

	gitdir := strings.Split(line[len(GitDirPrefix):], "\n")[0]
	gitdir = strings.TrimSpace(gitdir)
	return gitdir, nil
}
