package csheet

import (
	"fmt"

	"github.com/atotto/clipboard"
)

func PrintEntries(csheetFile string, copyToClipboard bool, quiet bool) {
	entries := GetEntries(csheetFile)

	if !quiet {
		for _, entry := range entries.Entries {
			fmt.Println(entry.String())
		}
	}

	if copyToClipboard {
		clipboard.WriteAll(entries.String())
	}
}
