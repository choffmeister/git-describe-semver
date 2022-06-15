package internal

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GenerateVersionOptions ...
type GenerateVersionOptions struct {
	FallbackTagName       string
	DropTagNamePrefix     bool
	PrereleaseSuffix      string
	PrereleasePrefix      string
	PrereleaseTimestamped bool
	Format                string
}

// GenerateVersion ...
func GenerateVersion(tagName string, counter int, headHash string, timestamp time.Time, opts GenerateVersionOptions) (*string, error) {
	devPrerelease := []string{opts.PrereleasePrefix, strconv.Itoa(counter), "g" + (headHash)[0:7]}
	if opts.PrereleaseTimestamped {
		timestampUTC := timestamp.UTC()
		timestampSegments := []string{
			strconv.FormatInt(timestampUTC.UnixMilli()/1000, 10),
		}
		devPrerelease = []string{opts.PrereleasePrefix, strings.Join(timestampSegments, ""), "g" + (headHash)[0:7]}
	}
	if opts.PrereleaseSuffix != "" {
		devPrerelease[len(devPrerelease)-1] = devPrerelease[len(devPrerelease)-1] + "-" + opts.PrereleaseSuffix
	}
	version := &SemVer{}
	if tagName == "" {
		version = SemVerParse(opts.FallbackTagName)
		if version == nil {
			return nil, fmt.Errorf("unable to parse fallback tag")
		}
		version.Prerelease = devPrerelease
	} else {
		version = SemVerParse(tagName)
		if version == nil {
			return nil, fmt.Errorf("unable to parse tag")
		}
		if counter > 0 {
			if len(version.Prerelease) > 0 {
				version = &SemVer{
					Prefix:        version.Prefix,
					Major:         version.Major,
					Minor:         version.Minor,
					Patch:         version.Patch,
					Prerelease:    append(version.Prerelease, devPrerelease...),
					BuildMetadata: append([]string{}, version.BuildMetadata...),
				}
			} else {
				version = &SemVer{
					Prefix:        version.Prefix,
					Major:         version.Major,
					Minor:         version.Minor,
					Patch:         version.Patch + 1,
					Prerelease:    devPrerelease,
					BuildMetadata: append([]string{}, version.BuildMetadata...),
				}
			}
		}
	}
	if opts.DropTagNamePrefix {
		version.Prefix = ""
	}
	result := version.String()
	if opts.Format != "" {
		result = strings.ReplaceAll(opts.Format, "<version>", result)
	}
	return &result, nil
}
