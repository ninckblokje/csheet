# csheet

This is a small app written in Go (my first Go app) for reading code cheat sheets from a Markdown document. `csheet` is [licensed](LICENSE) under the BSD-2-Clause.

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

## Other options

You can specify the file using `-f`

````
$ csheet -f csheet.md subject section
Stuff to remember
````

Results can be copied directory to the clipboard using `-c`:

````
$ csheet -c subject section
Stuff to remember
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
Usage of csheet:
  -c    Copy result to clipboard
  -e    Open editor using $EDITOR
  -f string
        Cheat sheet Mardown file
  -l    Show all possible entries
  -v    Display version
````

The options `f` can be combined with both `-e` or `-c`.
