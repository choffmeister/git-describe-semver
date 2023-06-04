package cmd

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/choffmeister/git-describe-semver/internal"
	"github.com/go-git/go-git/v5"
	"github.com/jessevdk/go-flags"
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

type ParserOptions struct {
	Dir string `long:"dir" default:"." description:"The git worktree directory"`
	Fallback string `long:"fallback" description:"The first version to fallback to should there be no tag"`
	DropPrefix bool `long:"drop-prefix" description:"Drop prefix from output"`
	PrereleaseSuffix string `long:"prerelease-suffix" description:"Suffix to add to prereleases"`
	PrereleasePrefix string `long:"prerelease-prefix" default:"dev" description:"Prefix to use as start of prerelease"`
	PrereleaseTimestamped bool `long:"prerelease-timestamped" description:"Use timestamp instead of commit count for prerelease"`
	NextRelease string `long:"next-release" description:"Bump current version to next release" choice:"major" choice:"minor" choice:"patch"`
	Format string `long:"format" description:"Format of output (use <version> as placeholder)"`
}

func Execute(version FullVersion) error {
	var options ParserOptions
	parser := flags.NewParser(&options, flags.Default)
	args, err := parser.Parse()
	if err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	opts := internal.GenerateVersionOptions{
		FallbackTagName:       options.Fallback,
		DropTagNamePrefix:     options.DropPrefix,
		PrereleaseSuffix:      options.PrereleaseSuffix,
		PrereleasePrefix:      options.PrereleasePrefix,
		PrereleaseTimestamped: options.PrereleaseTimestamped,
		NextRelease: 		   options.NextRelease,
		Format:                options.Format,
	}
	result, err := run(options.Dir, opts)
	if err != nil {
		return err
	}

	file := "-"
	if len(args) == 1 {
		arg := args[0]
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
