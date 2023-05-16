package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/choffmeister/git-describe-semver/internal"
	"github.com/go-git/go-git/v5"
)

func run(dir string, opts internal.GenerateVersionOptions) (*string, error) {
	repo, err := git.PlainOpen(dir)
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

func openStdoutOrFile(file string) (io.WriteCloser, error) {
	if file == "-" {
		return os.Stdout, nil
	}
	return os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
}

func Execute(version FullVersion) error {
	dirFlag := flag.String("dir", ".", "The git worktree directory")
	fallbackFlag := flag.String("fallback", "", "The first version to fallback to should there be no tag")
	dropPrefixFlag := flag.Bool("drop-prefix", false, "Drop prefix from output")
	prereleaseSuffixFlag := flag.String("prerelease-suffix", "", "Suffix to add to prereleases")
	prereleasePrefixFlag := flag.String("prerelease-prefix", "dev", "Prefix to use as start of prerelease (default to \"dev\"))")
	prereleaseTimestampedFlag := flag.Bool("prerelease-timestamped", false, "Use timestamp instead of commit count for prerelease")
	formatFlag := flag.String("format", "", "Format of output")
	flag.Parse()

	opts := internal.GenerateVersionOptions{
		FallbackTagName:       *fallbackFlag,
		DropTagNamePrefix:     *dropPrefixFlag,
		PrereleaseSuffix:      *prereleaseSuffixFlag,
		PrereleasePrefix:      *prereleasePrefixFlag,
		PrereleaseTimestamped: *prereleaseTimestampedFlag,
		Format:                *formatFlag,
	}
	result, err := run(*dirFlag, opts)
	if err != nil {
		return err
	}

	file := "-"
	if len(flag.Args()) == 1 {
		arg := flag.Args()[0]
		if strings.HasPrefix(arg, "$") {
			file = os.Getenv(strings.TrimPrefix(arg, "$"))
		} else {
			file = arg
		}
	}
	output, err := openStdoutOrFile(file)
	if err != nil {
		return err
	}
	defer output.Close()
	fmt.Fprintf(output, "%s\n", *result)

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
