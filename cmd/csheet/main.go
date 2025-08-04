package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ninckblokje/csheet/internal/csheet"
)

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

	var csheetFile string
	if *fileArg == "" {
		csheetFile = csheet.GetCSheetDir() + string(os.PathSeparator) + "csheet.md"
	} else {
		csheetFile = *fileArg
	}

	if *versionArg {
		printVersion()
	} else if *listArg {
		csheet.PrintEntries(csheetFile, *clipboardArg, *quietArg)
	} else {
		subject := args[0]
		section := args[1]

		csheet.PrintEntry(csheetFile, subject, section, *clipboardArg, *quietArg)
	}
}

func printUsage() {
	fmt.Println("Usage: csheet { OPTIONS } [SUBJECT] [SECTION]")
	fmt.Println("Options:")
	fmt.Println("-c        : Copy result to clipboard")
	fmt.Println("-f [FILE] : Specifies the Markdown file to read")
	fmt.Println("-h        : Print help")
	fmt.Println("-l        : Show all possible entries")
	fmt.Println("-q        : No output, useful with -c")
	fmt.Println("-v        : Shows the versions")
}

func printVersion() {
	versionInfo := strings.Split(version, "-")
	csheetVersion := strings.Join(versionInfo[:len(versionInfo)-1], "-")
	csheetRevision := versionInfo[len(versionInfo)-1]

	fmt.Printf("csheet version v%s, revision %s", csheetVersion, csheetRevision)
	fmt.Println("")
	fmt.Println("See: https://github.com/ninckblokje/csheet")
	fmt.Println("")
	fmt.Println("For my kids, L&M")
}

func validateArgs(fileArg *string, listArg *bool, versionArg *bool, args []string) {
	if *versionArg {
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
