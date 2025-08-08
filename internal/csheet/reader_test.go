package csheet

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestGetEntries(t *testing.T) {
	csheetFile := filepath.Join("..", "..", "csheet.md")
	entries := GetEntries(csheetFile)
	if len(entries.Entries) != 4 {
		t.Errorf("Expected 4 entries, got %d", len(entries.Entries))
	}

	expectedEntries := Entries{
		Entries: []Entry{
			{Subject: "TestSubjectOne", Section: "TestEntryTwo"},
			{Subject: "TestSubjectOne", Section: "TestEntryThree"},
			{Subject: "TestSubjectFour", Section: "TestEntryFive"},
			{Subject: "TestSubjectFour", Section: "TestEntrySix"},
		},
	}
	if !reflect.DeepEqual(expectedEntries, entries) {
		t.Errorf("Expected %s to be equal to %s", entries, expectedEntries)
	}
}
