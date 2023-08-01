package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/atotto/clipboard"
)

var csheetFile string
var csheetRevision = "DEV-BUILD"
var csheetVersion = "DEV-BUILD"

func main() {
	var clipboardArg = flag.Bool("c", false, "Copy result to clipboard")
	var editorArg = flag.Bool("e", false, "Open editor using $EDITOR")
	var fileArg = flag.String("f", "", "Cheat sheet Mardown file")
	var listArg = flag.Bool("l", false, "Show all possible entries")
	var quietArg = flag.Bool("q", false, "No output")
	var versionArg = flag.Bool("v", false, "Display version")

	flag.Parse()

	args := flag.Args()
	validateArgs(editorArg, fileArg, listArg, versionArg, args)

	if *fileArg == "" {
		csheetFile = getCSheetDir() + string(os.PathSeparator) + "csheet.md"
	} else {
		csheetFile = *fileArg
	}

	if *versionArg {
		printVersion()
	} else if *editorArg {
		openEditor()
	} else if *listArg {
		printEntries(*clipboardArg, *quietArg)
	} else {
		subject := args[0]
		section := args[1]

		printEntry(subject, section, *clipboardArg, *quietArg)
	}
}

func findEntry(fp *os.File, subject string, section string) []string {
	r := bufio.NewReaderSize(fp, 4*1024)
	demarcation := "## "
	if findHeader(r, "## "+subject, nil) && findHeader(r, "### "+section, &demarcation) {
		return readCode(r)
	}

	return nil
}

func findEntries(fp *os.File) []string {
	var entries []string
	var subject *string

	r := bufio.NewReaderSize(fp, 4*1024)
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

func getCSheetDir() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	return usr.HomeDir
}

func openEditor() {
	var editor = os.Getenv("EDITOR")
	if editor == "" {
		panic("$EDITOR is not defined")
	}

	cmd := exec.Command(editor, csheetFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(editor + " " + csheetFile)
		panic(err)
	}
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

func printEntry(subject string, section string, copyToClipboard bool, quiet bool) {
	fp := openFile()
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

func printEntries(copyToClipboard bool, quiet bool) {
	fp := openFile()
	defer fp.Close()

	entries := findEntries(fp)

	if !quiet {
		for i := 0; i < len(entries); i++ {
			fmt.Println(entries[i])
		}
	}

	if copyToClipboard {
		clipboardEntries := strings.Join(entries, "\n")
		clipboard.WriteAll(clipboardEntries)
	}
}

func printUsage() {
	fmt.Println("Usage: csheet { OPTIONS } [SUBJECT] [SECTION]")
	fmt.Println("Options:")
	fmt.Println("-c        : Copy result to clipboard")
	fmt.Println("-e        : Open editor using $EDITOR")
	fmt.Println("-f [FILE] : Specifies the Markdown file to read")
	fmt.Println("-h        : Print help")
	fmt.Println("-l        : Show all possible entries")
	fmt.Println("-q        : No output, useful with -c")
	fmt.Println("-v        : Shows the versions")
}

func printVersion() {
	fmt.Printf("csheet version %s, revision %s", csheetVersion, csheetRevision)
	fmt.Println("")
	fmt.Println("See: https://github.com/ninckblokje/csheet")
	fmt.Println("")
	fmt.Println("For my kids, L&M")
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

func validateArgs(editorArg *bool, fileArg *string, listArg *bool, versionArg *bool, args []string) {
	if *versionArg {
		// ok
		return
	} else if *editorArg && os.Getenv("EDITOR") != "" {
		// ok
		return
	} else if *listArg {
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
