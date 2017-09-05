# Folder tree generator from XLS data

## What this is

This small tool will generate a folder tree structure from XLS columns. For example

| MAINFOLDER    | Column2        | Column3        | Column4       |
|---------------|----------------|----------------|---------------|
| Main folder 1 |                |                |               |
|               | Second level 1 |                |               |
|               |                | Third level  1 |               |
|               |                | Third level 2  |               |
|               |                |                | Forth level 1 |
|               |                |                | Forth level 2 |
|               |                |                | Forth level 3 |
|               |                | Third level 3  |               |
| Main folder 2 |                |                |               |
| Main folder 3 |                |                |               |


## Configuration

Check `config.json` to change the xls path, generated path and the depth of folders to generate


## Run on windows

Run `launch.bat`


## Build

This is a `golang` application. If you want to build in your OS please check [Official golang website](https://golang.org/)




