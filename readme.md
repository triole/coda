# Coda ![build](https://github.com/triole/coda/actions/workflows/build.yaml/badge.svg)

<!-- toc -->

- [Synopsis](#synopsis)
- [Configuration](#configuration)

<!-- /toc -->

## Synopsis

...brings the music back to coding by running your linter, format and whatever scripts with a single command. It detects the kind of file that was passed and executes the appropriate commands.

## Configuration

All definitions are made in the `conf.toml`, which has to be in the same path as the executable binary. This is how a typical entry looks:

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

The `conf` folder in the repository provides a [full example](https://github.com/triole/coda/blob/master/conf/coda.toml).

As you can see the configuration is a list of entries. They function like this

| entry   | function                                                          |
|---------|-------------------------------------------------------------------|
| name    | only test relevant, provides a clearer output of the test results |
| regex   | 1st method file type detection                                    |
| shebang | 2nd method file type detection                                    |
| cmds    | list of commands to run                                           |

There is a set of variables that can be used inside a config. They get replaced by their appropriate values. This is what they look like:

```go mdox-exec="r -p"

Available variables

 VARIABLE              | DESCRIPTION                                 
-----------------------+---------------------------------------------
 {{.extension}}        | file's extension                            
 {{.filename}}         | full file name                              
 {{.filename_no_ext}}  | full file name without preceeding extension 
 {{.shortname}}        | short name, file name without path          
 {{.shortname_no_ext}} | short name without extension                

```

## How to use?

```shell
coda -h
coda eggs.py
coda main.rs

# debug mode, to just print commands that would have been run
coda -d shell_script
```
