package csheet

import (
	"reflect"
	"strings"
	"testing"
)

var testSheet = strings.Join([]string{
	"# csheet",
	"",
	"## go",
	"",
	"### build",
	"",
	"````",
	"go build ./...",
	"````",
	"",
	"### test",
	"",
	"````",
	"go test ./...",
	"````",
	"",
	"## git",
	"",
	"### status",
	"",
	"````",
	"git status",
	"````",
}, "\n")

func TestFindEntry(t *testing.T) {
	got := FindEntry(strings.NewReader(testSheet), "go", "build")
	want := []string{"go build ./..."}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FindEntry() = %v, want %v", got, want)
	}
}

func TestFindEntryMissingSection(t *testing.T) {
	got := FindEntry(strings.NewReader(testSheet), "go", "missing")
	if len(got) != 0 {
		t.Fatalf("FindEntry() = %v, want no lines", got)
	}
}

func TestFindEntries(t *testing.T) {
	got := FindEntries(strings.NewReader(testSheet))
	want := []string{"go build", "go test", "git status"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FindEntries() = %v, want %v", got, want)
	}
}

func TestFilterEntries(t *testing.T) {
	entries := []string{"go build", "go test", "git status"}
	got := FilterEntries(entries, "go")
	want := []string{"go build", "go test"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FilterEntries() = %v, want %v", got, want)
	}
}
