package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/choffmeister/git-describe-semver/internal"
	"github.com/go-git/go-git/v5"
)

// Run ...
func Run(dir string, opts internal.GenerateVersionOptions) (*string, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to open git repository: %v", err)
	}
	tagName, counter, headHash, err := internal.GitDescribe(*repo)
	if err != nil {
		return nil, fmt.Errorf("unable to describe commit: %v", err)
	}
	result, err := internal.GenerateVersion(*tagName, *counter, *headHash, opts)
	if err != nil {
		return nil, fmt.Errorf("unable to generate version: %v", err)
	}
	return result, nil
}

func main() {
	fallback := flag.String("fallback", "", "The first version to fallback to should there be no tag")
	dropPrefix := flag.Bool("drop-prefix", false, "Drop prefix from output")
	prereleaseSuffix := flag.String("prerelease-suffix", "", "Suffix to add to prereleases")
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to determine current directory: %v\n", err)
	}
	opts := internal.GenerateVersionOptions{
		FallbackTagName:   *fallback,
		DropTagNamePrefix: *dropPrefix,
		PrereleaseSuffix:  *prereleaseSuffix,
	}
	result, err := Run(dir, opts)
	if err != nil {
		log.Fatalf("unable to generate version: %v\n", err)
	}
	fmt.Println(*result)
}
