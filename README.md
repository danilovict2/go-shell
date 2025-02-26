# Go Shell

A simple implementation of a Unix shell in the Go programming language.

## Motivation

Modern shells like Bash have made everything so much easier. That's why I decided to get my hands dirty and use Go to explore how shells parse commands, execute programs, provide autocomplete functionality, and more.

## 🚀 Quick Start

Ensure you have a [Go](https://golang.org/doc/install) environment set up.

### Clone the project:

```bash
git https://github.com/danilovict2/go-shell.git
cd go-shell
```
### Run:

```bash
./your_program.sh
```

## 📖 Usage

### Builtin commands

* `exit`
* `echo`
* `type`
* `pwd`
* `cd`

## Examples

### Running a Program

To run a program that is available in your `PATH` (such as `ls`), simply type the program's name and press Enter:

```bash
ls
```

This will execute the `ls` command.

### Quoting

To handle spaces and special characters in arguments, you can use single and double quotes, as well as `\`. For example:

```bash
echo 'shell\"examplescript\"hello'
```

This will print:

```
shell\"examplescript\"hello
```

### Autocomplete

The shell provides autocomplete functionality for builtins and programs. To use autocomplete, start typing a command and press `Tab`. If there are multiple matches, press `Tab` again to see a list of possible completions.

For example, to autocomplete the `xyz_foo_bar_baz` command, you can type:

```bash
xyz_
```

Then press `Tab`, and the shell will complete it to:

```bash
xyz_foo_bar_baz
```

If there are multiple matches, pressing `Tab` again will show a list of possible completions.

## 🤝 Contributing

### Build the project

```bash
go build -o shell cmd/myshell/*.go
```

### Run the project

```bash
./shell
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `master` branch.