package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemVerString(t *testing.T) {
	assert := assert.New(t)
	test := func(input SemVer, expected string) {
		actual := input.String()
		assert.Equal(expected, actual)
	}

	test(SemVer{}, "0.0.0")
	test(SemVer{Prefix: "v"}, "v0.0.0")
	test(SemVer{Major: 1, Minor: 2, Patch: 3}, "1.2.3")
	test(SemVer{Prerelease: []string{"rc", "1"}}, "0.0.0-rc.1")
	test(SemVer{Prerelease: []string{"alpha-version", "1"}}, "0.0.0-alpha-version.1")
	test(SemVer{BuildMetadata: []string{"foo", "bar"}}, "0.0.0+foo.bar")
	test(SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3, Prerelease: []string{"rc", "1"}, BuildMetadata: []string{"foo", "bar"}}, "v1.2.3-rc.1+foo.bar")
}

func TestSemVerParse(t *testing.T) {
	assert := assert.New(t)
	test := func(input string, expected *SemVer) {
		actual := SemVerParse(input)
		assert.Equal(expected, actual)
	}

	test("0.0.0", &SemVer{})
	test("v0.0.0", &SemVer{Prefix: "v"})
	test("1.2.3", &SemVer{Major: 1, Minor: 2, Patch: 3})
	test("0.0.0-rc.1", &SemVer{Prerelease: []string{"rc", "1"}})
	test("0.0.0-alpha-version.1", &SemVer{Prerelease: []string{"alpha-version", "1"}})
	test("0.0.0+foo.bar", &SemVer{BuildMetadata: []string{"foo", "bar"}})
	test("v1.2.3-rc.1+foo.bar", &SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3, Prerelease: []string{"rc", "1"}, BuildMetadata: []string{"foo", "bar"}})
	test("invalid", nil)
}

func TestSemVerEqual(t *testing.T) {
	assert := assert.New(t)
	test := func(a SemVer, b SemVer, expected bool) {
		actual := a.Equal(b)
		assert.Equal(expected, actual)
	}

	test(SemVer{}, SemVer{}, true)
	test(SemVer{Major: 1}, SemVer{Major: 2}, false)
	test(SemVer{Minor: 1}, SemVer{Minor: 2}, false)
	test(SemVer{Patch: 1}, SemVer{Patch: 2}, false)
	test(SemVer{Prerelease: []string{"foo"}}, SemVer{Prerelease: []string{"foo"}}, true)
	test(SemVer{Prerelease: []string{"foo"}}, SemVer{Prerelease: []string{"bar"}}, false)
	test(SemVer{Prerelease: []string{"foo"}}, SemVer{}, false)
	test(SemVer{}, SemVer{Prerelease: []string{"bar"}}, false)
	test(SemVer{BuildMetadata: []string{"foo"}}, SemVer{BuildMetadata: []string{"foo"}}, true)
	test(SemVer{BuildMetadata: []string{"foo"}}, SemVer{BuildMetadata: []string{"bar"}}, false)
	test(SemVer{BuildMetadata: []string{"foo"}}, SemVer{}, false)
	test(SemVer{}, SemVer{BuildMetadata: []string{"bar"}}, false)
}
