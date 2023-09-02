# Coda ![build](https://github.com/triole/coda/actions/workflows/build.yaml/badge.svg)

<!-- toc -->

- [Synopsis](#synopsis)
- [Configuration](#configuration)
- [```go mdox-exec=&quot;cat examples/conf.yaml&quot;](#go-mdox-execcat-examplesconfyaml)
- [How to use?](#how-to-use)
- [Help](#help)

<!-- /toc -->

## Synopsis

...brings the music back to coding by running your linter, format and whatever scripts with a single command. It detects the kind of file that was passed and executes the appropriate commands.

## Configuration

Coda's config location can be set by using `-c`. If the flag is not provided coda tries to lookup the file's location. The following paths are tried in order. First match will be taken.

- coda binary folder + "coda.toml"
- ${HOME}/.conf/coda/conf.yaml
- ${HOME}/.conf/coda/conf.toml
- ${HOME}/.config/coda/conf.yaml
- ${HOME}/.config/coda/conf.toml

The configuration file can be toml or yaml and contains the filetype definitions and settings. A yaml example looks like this. Please look into [examples](https://github.com/triole/coda/blob/master/examples) for more information.

```go mdox-exec="cat examples/conf.yaml"
---
filetypes:
  - name: bash
    regex: ".*\\.sh$"
    shebang: "#!/bin/bash"
    cmds:
      - ["shellcheck", "-x", "-e", "SC1090,SC2154,SC2155", "{{.filename}}"]
      - ["shfmt", "-w", "-ci", "-i", "2", "{{.filename}}"]

  - name: golang
    regex: ".*\\.go$"
    cmds:
      - ["goimports", "-w", "{{.filename}}"]
      - ["gofmt", "-w", "{{.filename}}"]
      - ["staticcheck", "{{.filename}}"]
      - ["go", "test", "-v"]


settings:
  ignore_list:
    - ".git/"
    - ".pytest"
    - "__pycache__"
    - "/target/debug/"
    - "/target/release/"
    - ".lock$"
    - ".rs.racertmp"
```

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

## Help

```go mdox-exec="r -h"

brings the music back to coding

Arguments:
  [<filename>]    file to process, positional arg required

Flags:
  -h, --help            Show context-sensitive help.
  -c, --config="/home/ole/.conf/coda/conf.yaml"
                        configuration file
  -p, --print-vars      print available vars
  -n, --dry-run         dry run, just print don't do
  -V, --version-flag    display version
```
