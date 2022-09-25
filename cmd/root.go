package cmd

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/choffmeister/git-describe-semver/internal"
)

func run(dir string, opts internal.GenerateVersionOptions) (*string, error) {
	repo, err := internal.OpenRepository(dir)
	if err != nil {
		return nil, fmt.Errorf("unable to open git repository: %v", err)
	}
	tagName, counter, headHash, err := internal.GitDescribe(*repo)
	if err != nil {
		return nil, fmt.Errorf("unable to describe commit: %v", err)
	}
	result, err := internal.GenerateVersion(*tagName, *counter, *headHash, time.Now(), opts)
	if err != nil {
		return nil, fmt.Errorf("unable to generate version: %v", err)
	}
	return result, nil
}

func Execute(version FullVersion) error {
	fallbackFlag := flag.String("fallback", "", "The first version to fallback to should there be no tag")
	dropPrefixFlag := flag.Bool("drop-prefix", false, "Drop prefix from output")
	prereleaseSuffixFlag := flag.String("prerelease-suffix", "", "Suffix to add to prereleases")
	prereleasePrefixFlag := flag.String("prerelease-prefix", "dev", "Prefix to use as start of prerelease (default to \"dev\"))")
	prereleaseTimestampedFlag := flag.Bool("prerelease-timestamped", false, "Use timestamp instead of commit count for prerelease")
	formatFlag := flag.String("format", "", "Format of output")
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	opts := internal.GenerateVersionOptions{
		FallbackTagName:       *fallbackFlag,
		DropTagNamePrefix:     *dropPrefixFlag,
		PrereleaseSuffix:      *prereleaseSuffixFlag,
		PrereleasePrefix:      *prereleasePrefixFlag,
		PrereleaseTimestamped: *prereleaseTimestampedFlag,
		Format:                *formatFlag,
	}
	result, err := run(dir, opts)
	if err != nil {
		return err
	}
	fmt.Println(*result)

	return nil
}

type FullVersion struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

func (v FullVersion) ToString() string {
	result := v.Version
	if v.Commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, v.Commit)
	}
	if v.Date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, v.Date)
	}
	if v.BuiltBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, v.BuiltBy)
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		result = fmt.Sprintf("%s\nmodule version: %s, checksum: %s", result, info.Main.Version, info.Main.Sum)
	}
	return result
}
