# csheet

This is a small app written in Go (my first Go app) for reading code cheat sheets from a Markdown document.

It expects the file `csheet.md` to be in the users home directory. Or optionally it is possible to specify which file should be used. 

Cheat sheets follow the following structuur:

`````markdown
# csheet

## subject

### section

````
Stuff to remember
````

`````

Then retrieve it using this command (from `csheet.md` in the users home directory):
````
$ csheet subject section
Stuff to remember
````

Or specify the file manually:
````
$ csheet -f csheet.md subject section
Stuff to remember
````

To get all the subjects and sections, use `-l` (or combined with `-f`):
````
$ csheet -l
subject section
````

To get the version use `-v`.