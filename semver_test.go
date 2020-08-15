package main

import (
	"testing"
)

func TestSemVerString(t *testing.T) {
	test := func(input SemVer, expected string) {
		actual := input.String()
		if actual != expected {
			t.Errorf("expected %v to be string as %v, got %v", input, expected, actual)
		}
	}

	test(SemVer{}, "0.0.0")
	test(SemVer{Prefix: "v"}, "v0.0.0")
	test(SemVer{Major: 1, Minor: 2, Patch: 3}, "1.2.3")
	test(SemVer{PreRelease: []string{"rc", "1"}}, "0.0.0-rc.1")
	test(SemVer{PreRelease: []string{"alpha-version", "1"}}, "0.0.0-alpha-version.1")
	test(SemVer{BuildMetadata: []string{"foo", "bar"}}, "0.0.0+foo.bar")
	test(SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3, PreRelease: []string{"rc", "1"}, BuildMetadata: []string{"foo", "bar"}}, "v1.2.3-rc.1+foo.bar")
}

func TestSemVerParse(t *testing.T) {
	test := func(input string, expected *SemVer) {
		actual := SemVerParse(input)
		if actual == nil && expected == nil {
			// ok
		} else if actual == nil || expected == nil || !actual.Equal(*expected) {
			t.Errorf("expected %v to be parsed as %v, got %v", input, expected, actual)
		}
	}

	test("0.0.0", &SemVer{})
	test("v0.0.0", &SemVer{Prefix: "v"})
	test("1.2.3", &SemVer{Major: 1, Minor: 2, Patch: 3})
	test("0.0.0-rc.1", &SemVer{PreRelease: []string{"rc", "1"}})
	test("0.0.0-alpha-version.1", &SemVer{PreRelease: []string{"alpha-version", "1"}})
	test("0.0.0+foo.bar", &SemVer{BuildMetadata: []string{"foo", "bar"}})
	test("v1.2.3-rc.1+foo.bar", &SemVer{Prefix: "v", Major: 1, Minor: 2, Patch: 3, PreRelease: []string{"rc", "1"}, BuildMetadata: []string{"foo", "bar"}})
	test("invalid", nil)
}
