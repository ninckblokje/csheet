package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/ninckblokje/csheet/internal/csheet"
)

var csheetFile string
var commit = "test"
var version = "DEV-BUILD"

func main() {
	var clipboardArg = flag.Bool("c", false, "Copy result to clipboard")
	var fileArg = flag.String("f", "", "Cheat sheet Mardown file")
	var listArg = flag.Bool("l", false, "Show all possible entries")
	var quietArg = flag.Bool("q", false, "No output")
	var versionArg = flag.Bool("v", false, "Display version")

	flag.Parse()

	args := flag.Args()
	validateArgs(fileArg, listArg, versionArg, args)

	if *fileArg == "" {
		csheetFile = getCSheetDir() + string(os.PathSeparator) + "csheet.md"
	} else {
		csheetFile = *fileArg
	}

	if *versionArg {
		printVersion()
	} else if *listArg && len(args) == 0 {
		printEntries(*clipboardArg, *quietArg)
	} else if *listArg && len(args) == 1 {
		subject := args[0]

		printEntriesForSubject(*clipboardArg, *quietArg, subject)
	} else {
		subject := args[0]
		section := args[1]

		printEntry(subject, section, *clipboardArg, *quietArg)
	}
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

func printEntry(subject string, section string, copyToClipboard bool, quiet bool) {
	fp := openFile()
	defer fp.Close()

	code := csheet.FindEntry(fp, subject, section)

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

	entries := csheet.FindEntries(fp)

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

func printEntriesForSubject(copyToClipboard bool, quiet bool, subject string) {
	fp := openFile()
	defer fp.Close()

	entries := csheet.FilterEntries(csheet.FindEntries(fp), subject)

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
	fmt.Println("-c           : Copy result to clipboard")
	fmt.Println("-f [FILE]    : Specifies the Markdown file to read")
	fmt.Println("-h           : Print help")
	fmt.Println("-l {SUBJECT} : Show all possible entries")
	fmt.Println("-q           : No output, useful with -c")
	fmt.Println("-v           : Shows the versions")
}

func printVersion() {
	fmt.Printf("csheet version v%s, revision %s", version, commit)
	fmt.Println("")
	fmt.Println("See: https://github.com/ninckblokje/csheet")
	fmt.Println("")
	fmt.Println("For my kids, L&M")
}

func validateArgs(fileArg *string, listArg *bool, versionArg *bool, args []string) {
	if *versionArg {
		// ok
		return
	} else if *listArg && len(args) <= 1 {
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
