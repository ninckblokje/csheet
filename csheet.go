package main

import (
	"os"
	"bufio"
	"io"
	"strings"
	"fmt"
	"os/user"
	"flag"
)

var csheetFile string
var csheetVersion = "DEV-BUILD"

func main() {
	var fileArg = flag.String("f", "", "Cheat sheet Mardown file")
	var versionArg = flag.Bool("v", false, "Display version")

	flag.Parse()

	args := flag.Args()
	validateArgs(fileArg, versionArg, args)

	if *fileArg == "" {
		csheetFile = getCSheetDir() + string(os.PathSeparator) + "csheet.md"
	} else {
		csheetFile = *fileArg
	}

	if *versionArg {
		printVersion()
	} else {
		subject := args[0]
		section := args[1]

		readEntry(subject, section)
	}
}

func find(fp *os.File, subject string, section string) []string {
	r := bufio.NewReaderSize(fp, 4*1024)
	if findHeader(r, "## " + subject) && findHeader(r, "### " + section) {
		return readCode(r)
	}

	return nil
}

func findHeader(r *bufio.Reader, header string) bool {
	line := readLine(r)
	for line != nil {
		s := *line

		if strings.HasPrefix(s, header) {
			return true
		}

		line = readLine(r)
	}

	return false
}

func getCSheetDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return usr.HomeDir
}

func openFile() (fp *os.File) {
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

func printUsage() {
	fmt.Println("Usage: csheet { OPTIONS } [SUBJECT] [SECTION]")
	fmt.Println("Options:")
	fmt.Println("-f [FILE] : Specifies the Markdown file to read")
	fmt.Println("-v        : Shows the versions")
}

func printVersion() {
	fmt.Printf("csheet version %s", csheetVersion)
	fmt.Println("")
	fmt.Println("See: https://github.com/ninckblokje/csheet")
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

func readEntry(subject string, section string) {
	fp := openFile()
	defer fp.Close()

	code := find(fp, subject, section)
	for i := 0; i< len(code); i++ {
		fmt.Println(code[i])
	}
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

func validateArgs(fileArgs *string, versionArg *bool, args []string) {
	if *versionArg {
		// ok
		return
	} else if len(args) == 2 {
		// ok
		return
	} else {
		printUsage()
		os.Exit(1)
	}
}

func writeHeader(fp *os.File, header string) {
	w := bufio.NewWriter(fp)

	fmt.Fprintln(w, header)

	w.Flush()
}