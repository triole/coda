# Coda

<!--- mdtoc: toc begin -->

1.	[Synopsis](#synopsis)
2.	[Configuration](#configuration)
3.	[How to use?](#how-to-use-)<!--- mdtoc: toc end -->

## Synopsis

...brings the music back to coding by running your linter, format and whatever scripts with a single command. It detects the kind of file that was passed and executes the appropriate commands.

## Configuration

All definitions are made in the `conf.toml`. This is how a typical entry looks:

```toml
[[ft]]
name = "python"
regex = ".*\\.py$"
shebang = "#!/usr/bin/python3"
cmds = [
    ["autoflake", "--remove-all-unused-imports", "-i", "{FILENAME}"],
    ["isort", "{FILENAME}"],
    ["pytest", "-v", "{SHORTNAME_NO_EXT}_test.py"]
]
```

As you can see the configuration is a list of entries. They function like this

| entry   | function                                                          |
|---------|-------------------------------------------------------------------|
| name    | only test relevant, provides a clearer output of the test results |
| regex   | 1st method file type detection                                    |
| shebang | 2nd method file type detection                                    |
| cmds    | list of commands to run                                           |

There is a set of variables available. If used in command they are replaced by a part if the input file name. I hope they are be quite self explanatory. This is what they look like:

| var                | explanation                              |
|--------------------|------------------------------------------|
| {FILENAME}         | full file name including path            |
| {FILENAME_NO_EXT}  | like the one above but without extension |
| {SHORTNAME}        | file name without path                   |
| {SHORTNAME_NO_EXT} | above, but missing extension             |
| {EXT}              | file extension only                      |

## How to use?

```shell
coda -h
coda eggs.py
coda main.rs

# debug mode, to just print commands that would have been run
coda -d shell_script
```
