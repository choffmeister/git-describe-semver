package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateVersion(t *testing.T) {
	now, _ := time.Parse(time.RFC822Z, "01 Jan 20 02:03 -0000")
	assert := assert.New(t)
	test := func(inputTagName string, inputCounter int, inputHeadHash string, inputOpts GenerateVersionOptions, expected string) {
		actual, err := GenerateVersion(inputTagName, inputCounter, inputHeadHash, now, inputOpts)
		if assert.NoError(err) {
			assert.Equal(expected, *actual)
		}
	}

	test("0.0.0", 0, "abc1234", GenerateVersionOptions{PrereleasePrefix: "dev"}, "0.0.0")
	test("0.0.0", 1, "abc1234", GenerateVersionOptions{PrereleasePrefix: "dev"}, "0.0.1-dev.1.gabc1234")
	test("0.0.0-rc1", 1, "abc1234", GenerateVersionOptions{PrereleasePrefix: "dev"}, "0.0.0-rc1.dev.1.gabc1234")
	test("0.0.0-rc.1", 1, "abc1234", GenerateVersionOptions{PrereleasePrefix: "dev"}, "0.0.0-rc.1.dev.1.gabc1234")
	test("0.0.0-rc.1+foobar", 1, "abc1234", GenerateVersionOptions{PrereleasePrefix: "dev"}, "0.0.0-rc.1.dev.1.gabc1234+foobar")
	test("v0.0.0-rc.1+foobar", 1, "abc1234", GenerateVersionOptions{PrereleasePrefix: "dev"}, "v0.0.0-rc.1.dev.1.gabc1234+foobar")

	test("", 1, "abc1234", GenerateVersionOptions{FallbackTagName: "0.0.0", PrereleasePrefix: "dev"}, "0.0.0-dev.1.gabc1234")
	test("", 1, "abc1234", GenerateVersionOptions{FallbackTagName: "v0.0.0", PrereleasePrefix: "dev"}, "v0.0.0-dev.1.gabc1234")

	test("v0.0.0", 0, "abc1234", GenerateVersionOptions{PrereleaseSuffix: "SNAPSHOT", PrereleasePrefix: "dev"}, "v0.0.0")
	test("v0.0.0", 1, "abc1234", GenerateVersionOptions{PrereleaseSuffix: "SNAPSHOT", PrereleasePrefix: "dev"}, "v0.0.1-dev.1.gabc1234-SNAPSHOT")

	test("v0.0.0", 0, "abc1234", GenerateVersionOptions{DropTagNamePrefix: true, PrereleasePrefix: "dev"}, "0.0.0")
	test("v0.0.0-rc.1", 1, "abc1234", GenerateVersionOptions{DropTagNamePrefix: true, PrereleasePrefix: "dev"}, "0.0.0-rc.1.dev.1.gabc1234")
	test("v0.0.0-rc.1+foobar", 1, "abc1234", GenerateVersionOptions{DropTagNamePrefix: true, PrereleasePrefix: "dev"}, "0.0.0-rc.1.dev.1.gabc1234+foobar")
	test("", 1, "abc1234", GenerateVersionOptions{FallbackTagName: "v0.0.0", DropTagNamePrefix: true, PrereleasePrefix: "dev"}, "0.0.0-dev.1.gabc1234")

	test("0.0.0", 0, "abc1234", GenerateVersionOptions{PrereleasePrefix: "custom"}, "0.0.0")
	test("0.0.0", 1, "abc1234", GenerateVersionOptions{PrereleasePrefix: "custom"}, "0.0.1-custom.1.gabc1234")

	test("0.0.0", 0, "abc1234", GenerateVersionOptions{PrereleasePrefix: "dev", PrereleaseTimestamped: true}, "0.0.0")
	test("0.0.0", 1, "abc1234", GenerateVersionOptions{PrereleasePrefix: "dev", PrereleaseTimestamped: false}, "0.0.1-dev.1.gabc1234")
	test("0.0.0", 1, "abc1234", GenerateVersionOptions{PrereleasePrefix: "dev", PrereleaseTimestamped: true}, "0.0.1-dev.1577844180.gabc1234")

	_, err := GenerateVersion("", 1, "abc1234", now, GenerateVersionOptions{PrereleasePrefix: "dev"})
	assert.Error(err)
}
