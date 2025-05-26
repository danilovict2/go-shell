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
* `history`

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

### Pipelines

A pipeline connects the standard output of one command to the standard input of the next command using the | operator. For example:

```bash
ls /tmp/foo/file | wc
```

This will print:

```bash
      7       7      70
```

Pipelining more than two commands is possible. For example:
```bash
ls -la /tmp | tail -n 5 | head -n 3 | grep "file"
```

Will print:

```bash
-rw-rw-r--  1 user user       0 May  7 18:57 file
```

## Command History

The shell supports command history navigation using the arrow keys. You can recall and reuse previous commands as follows:

- Press <kbd>↑</kbd> (Up Arrow) to cycle backward through your command history.
- Press <kbd>↓</kbd> (Down Arrow) to move forward in the history.

Example session:

```bash
$ echo hello
hello
$ echo world
world
# Press ↑
$ echo world
# Press ↑ again
$ echo hello
# Press ↓
$ echo world
# Press Enter
world
$
```

You can also view and interact with your command history using the built-in `history` command.

```bash
$ history
1  echo hello
2  echo world
3  ls -l
```

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