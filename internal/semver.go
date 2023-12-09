package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// SemVer ...
type SemVer struct {
	Prefix        string
	Major         int
	Minor         int
	Patch         int
	Prerelease    []string
	BuildMetadata []string
}

// Equal ...
func (v SemVer) Equal(v2 SemVer) bool {
	return true &&
		v.Prefix == v2.Prefix &&
		v.Major == v2.Major &&
		v.Minor == v2.Minor &&
		v.Patch == v2.Patch &&
		equalStringSlice(v.Prerelease, v2.Prerelease) &&
		equalStringSlice(v.BuildMetadata, v2.BuildMetadata)
}

// Bump ...
func (v *SemVer) Bump(nextRelease string) {
	if nextRelease == "" {
		return
	}
	isPrerelease := len(v.Prerelease) > 0
	patch := v.Patch
	if nextRelease == "patch" && !isPrerelease {
		patch++
	}
	minor := v.Minor
	if nextRelease == "minor" {
		if v.Patch != 0 || !isPrerelease {
			minor++
		}
		patch = 0
	}
	major := v.Major
	if nextRelease == "major" {
		if v.Patch != 0 || !isPrerelease {
			major++
		}
		minor = 0
		patch = 0
	}
	v.Major = major
	v.Minor = minor
	v.Patch = patch
	v.Prerelease = []string{}
	v.BuildMetadata = []string{}
}

// String ...
func (v SemVer) String() string {
	str := fmt.Sprintf("%s%d.%d.%d", v.Prefix, v.Major, v.Minor, v.Patch)
	if len(v.Prerelease) > 0 {
		str = str + "-" + strings.Join(v.Prerelease, ".")
	}
	if len(v.BuildMetadata) > 0 {
		str = str + "+" + strings.Join(v.BuildMetadata, ".")
	}
	return str
}

var semVerRegexp = regexp.MustCompile(`^([A-Za-z]+)?(\d+)\.(\d+)\.(\d+)(?:-((?:[0-9A-Za-z-]+)(?:\.[0-9A-Za-z-]+)*))?(?:\+((?:[0-9A-Za-z-]+)(?:\.[0-9A-Za-z-]+)*))?$`)

// SemVerParse ...
func SemVerParse(str string) *SemVer {
	match := semVerRegexp.FindStringSubmatch(str)
	if len(match) == 0 {
		return nil
	}

	prefix := match[1]
	major, _ := strconv.Atoi(match[2])
	minor, _ := strconv.Atoi(match[3])
	patch, _ := strconv.Atoi(match[4])
	prerelease := stringToSlice(match[5], ".")
	buildMetadata := stringToSlice(match[6], ".")

	return &SemVer{
		Prefix:        prefix,
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		Prerelease:    prerelease,
		BuildMetadata: buildMetadata,
	}
}

func stringToSlice(s string, sep string) []string {
	temp := strings.Split(s, sep)
	if temp[0] == "" {
		return []string(nil)
	}
	return temp
}

func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
