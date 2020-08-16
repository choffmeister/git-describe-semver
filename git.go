package main

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitTagMap ...
func GitTagMap(repo git.Repository) (*map[string]string, error) {
	iter, err := repo.Tags()
	if err != nil {
		return nil, err
	}
	tagMap := map[string]string{}
	iter.ForEach(func(r *plumbing.Reference) error {
		tagMap[r.Hash().String()] = r.Name().Short()
		return nil
	})
	return &tagMap, nil
}

// GitDescribe ...
func GitDescribe(repo git.Repository) (*string, *int, *string, error) {
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
		Order: git.LogOrderDefault,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to get log: %v", err)
	}
	counter := 0
	tagHash := ""
	commits.ForEach(func(c *object.Commit) error {
		_, found := (*tags)[c.Hash.String()]
		if tagHash == "" && found {
			tagHash = c.Hash.String()
		}
		if tagHash == "" {
			counter = counter + 1
		}
		return nil
	})
	if tagHash == "" {
		tagName := ""
		return &tagName, &counter, &headHash, nil
	}
	tagName := (*tags)[tagHash]
	return &tagName, &counter, &headHash, nil
}
