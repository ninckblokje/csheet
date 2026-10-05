# csheet

This is a small app written in Go (my first Go app) for reading code cheat sheets from a Markdown document. `csheet` is [licensed](LICENSE) under the BSD-2-Clause.

## Build locally

Run `make` to build both command-line tools. The binaries are written to `bin/csheet` and `bin/csheet-mcp`.

## MCP server

`csheet-mcp` exposes cheat sheet entries to MCP clients over standard input and output. Add a server entry to your MCP client configuration, replacing the example paths with absolute paths on your machine:

```json
{
  "mcpServers": {
    "csheet": {
      "command": "/absolute/path/to/csheet/bin/csheet-mcp",
      "args": ["-f", "/absolute/path/to/csheet.md"]
    }
  }
}
```

The `-f` argument is optional. Without it, the server reads `csheet.md` from your home directory.

## Cheat sheet

By default `csheet.md` from the users home directory is read, but it is possible to specify a custom file using `-f`.

Cheat sheets follow the following structuur:

`````markdown
# csheet

## subject

### section

````
Stuff to remember
````

`````

`csheet` will print the content enclosed by four backticks.

## Get single section

Retrieve a single section:

````
$ csheet subject section
Stuff to remember
````

## List all subjects and sections

To get all the subjects and sections, use `-l`:

````
$ csheet -l
subject section
````

It is also possible to list only the sections for a specific subject. Use `-l {SUBJECT}` for that:

````
$ csheet -l subject
subject section
````

## Other options

You can specify the file using `-f`

````
$ csheet -f csheet.md subject section
Stuff to remember
````

Results can be copied directory to the clipboard using `-c` (install `xclip` for this to work):

````
$ csheet -c subject section
Stuff to remember
````

Or combined with `-q` output can be suppressed on the command line.

````
$ csheet -q -c subject section
````

The version can be printed using -v:

````
$ csheet -v
csheet version 1.4, revision 7da242a
See: https://github.com/ninckblokje/csheet
````

Help can be printed with `-h`:

````
$ csheet -h
Usage of ./csheet:
  -c    Copy result to clipboard
  -f string
        Cheat sheet Mardown file
  -l    Show all possible entries
  -q    No output
  -v    Display version
````

The options `f` can be combined with `-e`, `-q` or `-c`.
