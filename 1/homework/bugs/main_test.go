package main

import (
	"testing"
)

func TestShadowing(t *testing.T) {
	expected := 6142
	result := Shadowing()
	if result != expected {
		t.Errorf("results not match\nGot: %v\nExpected: %v", result, expected)
	}
}

func TestBadMap(t *testing.T) {
	err := BadMap()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestExistCounter(t *testing.T) {
	expected := 4
	result := ExistCounter()
	if result != expected {
		t.Errorf("results not match\nGot: %v\nExpected: %v", result, expected)
	}
}
