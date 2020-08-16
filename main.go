package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-git/go-git/v5"
)

// GenerateVersion ...
func GenerateVersion(tagName string, counter int, headHash string) (*string, error) {
	devPreRelease := []string{"dev", strconv.Itoa(counter)}
	buildMetadata := []string{"g" + (headHash)[0:7]}
	if tagName == "" {
		version := SemVer{
			Prefix:        "v",
			Major:         0,
			Minor:         0,
			Patch:         0,
			PreRelease:    devPreRelease,
			BuildMetadata: buildMetadata,
		}
		result := version.String()
		return &result, nil
	}
	version := SemVerParse(tagName)
	if version == nil {
		return nil, fmt.Errorf("unable to parse tag")
	}
	if counter == 0 {
		result := version.String()
		return &result, nil
	}
	if len(version.PreRelease) > 0 {
		version = &SemVer{
			Prefix:        version.Prefix,
			Major:         version.Major,
			Minor:         version.Minor,
			Patch:         version.Patch,
			PreRelease:    append(version.PreRelease, devPreRelease...),
			BuildMetadata: buildMetadata,
		}
	} else {
		version = &SemVer{
			Prefix:        version.Prefix,
			Major:         version.Major,
			Minor:         version.Minor,
			Patch:         version.Patch + 1,
			PreRelease:    devPreRelease,
			BuildMetadata: buildMetadata,
		}
	}
	result := version.String()
	return &result, nil
}

// Run ...
func Run(dir string) (*string, error) {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to open git repository: %v", err)
	}
	tagName, counter, headHash, err := GitDescribe(*repo)
	if err != nil {
		return nil, fmt.Errorf("unable to find head: %v", err)
	}
	result, err := GenerateVersion(*tagName, *counter, *headHash)
	if err != nil {
		return nil, fmt.Errorf("unable to generate version: %v", err)
	}
	return result, nil
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to determine current directory: %v\n", err)
	}
	result, err := Run(dir)
	if err != nil {
		log.Fatalf("unable to generate version: %v\n", err)
	}
	fmt.Println(*result)
}
