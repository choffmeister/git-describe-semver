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
	devPreRelease := []string{"dev", strconv.Itoa(counter), "g" + (headHash)[0:7]}
	if tagName == "" {
		version := SemVer{
			Prefix:     "v",
			Major:      0,
			Minor:      0,
			Patch:      0,
			PreRelease: devPreRelease,
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
			BuildMetadata: version.BuildMetadata,
		}
	} else {
		version = &SemVer{
			Prefix:        version.Prefix,
			Major:         version.Major,
			Minor:         version.Minor,
			Patch:         version.Patch + 1,
			PreRelease:    devPreRelease,
			BuildMetadata: version.BuildMetadata,
		}
	}
	result := version.String()
	return &result, nil
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to determine current directory: %v\n", err)
	}
	repo, err := git.PlainOpen(dir)
	if err != nil {
		log.Fatalf("unable to open git repository: %v\n", err)
	}
	tagName, counter, headHash, err := Describe(*repo)
	if err != nil {
		log.Fatalf("unable to find head: %v\n", err)
	}
	result, err := GenerateVersion(*tagName, *counter, *headHash)
	if err != nil {
		log.Fatalf("unable to generate version: %v\n", err)
	}
	fmt.Println(*result)
}