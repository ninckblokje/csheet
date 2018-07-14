# csheet

This is a small app written in Go (my first Go app) for reading code cheat sheets from a Markdown document.

It expects the file `csheet.md` to be in the users home directory. Cheat sheets follow the following structuur:

````markdown
# csheet

## subject

### section

```
Stuff to remember
```

````

Then retrieve it using this command:
````
$ csheet subject section
Stuff to remember
````