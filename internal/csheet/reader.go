package csheet

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/atotto/clipboard"
)

func findEntry(fp *os.File, subject string, section string) []string {
	r := bufio.NewReaderSize(fp, 4*1024)
	demarcation := "## "
	if findHeader(r, "## "+subject, nil) && findHeader(r, "### "+section, &demarcation) {
		return readCode(r)
	}

	return nil
}

func findEntries(fp *os.File) []Entry {
	var entries []Entry
	var subject *string

	r := bufio.NewReaderSize(fp, 4*1024)
	line := readLine(r)
	for line != nil {
		s := *line

		if strings.HasPrefix(s, "## ") {
			tmp := strings.TrimPrefix(s, "## ")
			subject = &tmp
		} else if strings.HasPrefix(s, "### ") && subject != nil {
			entries = append(entries, Entry{
				Subject: *subject,
				Section: strings.TrimPrefix(s, "### ")})
		}

		line = readLine(r)
	}

	log.Printf("Found %d entries", len(entries))
	return entries
}

func findHeader(r *bufio.Reader, header string, demarcation *string) bool {
	line := readLine(r)
	for line != nil {
		s := *line

		if s == header {
			return true
		} else if demarcation != nil && strings.HasPrefix(s, *demarcation) {
			return false
		}

		line = readLine(r)
	}

	return false
}

func GetCSheetDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return usr.HomeDir
}

func GetEntries(csheetFile string) Entries {
	fp := openFile(csheetFile)
	defer fp.Close()

	entries := findEntries(fp)
	return Entries{
		Entries: entries,
	}
}

func openFile(csheetFile string) (fp *os.File) {
	log.Printf("Opening CSheet file %s\n", csheetFile)

	_, err := os.Stat(csheetFile)
	if err == nil {
		fp, err := os.Open(csheetFile)
		if err != nil {
			panic(err)
		}
		return fp
	} else if os.IsNotExist(err) {
		fp, err := os.Create(csheetFile)
		if err != nil {
			panic(err)
		}

		writeHeader(fp, "# csheet")
		return fp
	} else {
		panic(err)
	}
}

func PrintEntryValue(csheetFile string, subject string, section string, copyToClipboard bool, quiet bool) {
	fp := openFile(csheetFile)
	defer fp.Close()

	code := findEntry(fp, subject, section)

	if !quiet {
		for i := 0; i < len(code); i++ {
			fmt.Println(code[i])
		}
	}

	if copyToClipboard {
		clipboardEntries := strings.Join(code, "\n")
		clipboard.WriteAll(clipboardEntries)
	}
}

func readCode(r *bufio.Reader) []string {
	var code []string
	var readCode = false

	line := readLine(r)
	for line != nil {
		s := *line

		if strings.HasPrefix(s, "````") {
			readCode = !readCode

			if !readCode {
				break
			}
		} else if readCode {
			code = append(code, s)
		}

		line = readLine(r)
	}

	return code
}

func readLine(r *bufio.Reader) *string {
	line, isPrefix, err := r.ReadLine()
	if isPrefix {
		panic("buffer size to small")
	}

	if err == nil {
		s := string(line)
		return &s
	} else if err != io.EOF {
		panic(err)
	}

	return nil
}

func writeHeader(fp *os.File, header string) {
	w := bufio.NewWriter(fp)

	fmt.Fprintln(w, header)

	w.Flush()
}
