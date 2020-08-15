package main

import (
	"testing"
)

func TestGenerateVersion(t *testing.T) {
	test := func(inputTagName string, inputCounter int, inputHeadHash string, expected string) {
		actual, err := GenerateVersion(inputTagName, inputCounter, inputHeadHash)
		if err != nil {
			t.Errorf("expected %v,%v,%v to be string as %v, got %v", inputTagName, inputCounter, inputHeadHash, expected, err)
		} else if *actual != expected {
			t.Errorf("expected %v,%v,%v to be string as %v, got %v", inputTagName, inputCounter, inputHeadHash, expected, actual)
		}
	}

	test("0.0.0", 0, "abc1234", "0.0.0")
	test("0.0.0", 1, "abc1234", "0.0.1-dev.1.gabc1234")
	test("0.0.0-rc1", 1, "abc1234", "0.0.0-rc1.dev.1.gabc1234")
	test("0.0.0-rc.1", 1, "abc1234", "0.0.0-rc.1.dev.1.gabc1234")
}
