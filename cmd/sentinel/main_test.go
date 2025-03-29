package main

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if version == "" {
		t.Error("Version should not be empty")
	}
	if commit == "" {
		t.Error("Commit should not be empty")
	}
	if date == "" {
		t.Error("Date should not be empty")
	}
}
