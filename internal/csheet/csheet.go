package csheet

import (
	"bufio"
	"io"
	"strings"
)

func FilterEntries(entries []string, subject string) []string {
	var filteredEntries []string

	for _, entry := range entries {
		if strings.HasPrefix(entry, subject+" ") {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	return filteredEntries
}

func FindEntry(src io.Reader, subject string, section string) []string {
	r := bufio.NewReaderSize(src, 4*1024)
	demarcation := "## "
	if findHeader(r, "## "+subject, nil) && findHeader(r, "### "+section, &demarcation) {
		return readCode(r)
	}

	return nil
}

func FindEntries(src io.Reader) []string {
	var entries []string
	var subject *string

	r := bufio.NewReaderSize(src, 4*1024)
	line := readLine(r)
	for line != nil {
		s := *line

		if strings.HasPrefix(s, "## ") {
			tmp := strings.TrimPrefix(s, "## ")
			subject = &tmp
		} else if strings.HasPrefix(s, "### ") && subject != nil {
			entries = append(entries, *subject+" "+strings.TrimPrefix(s, "### "))
		}

		line = readLine(r)
	}

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

func readCode(r *bufio.Reader) []string {
	var code []string
	var readingCode bool

	line := readLine(r)
	for line != nil {
		s := *line

		if strings.HasPrefix(s, "````") {
			readingCode = !readingCode

			if !readingCode {
				break
			}
		} else if readingCode {
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
