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

The shell supports command history navigation using the arrow keys, allowing you to recall and reuse previous commands. You can also view and interact with your history using the built-in `history` command.

**Note:** Command history is stored only for the current session and will be lost once the program exits.  
To make history persistent, use the file-based approach described below.

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

To view your command history, use:

```bash
$ history
1  echo hello
2  echo world
3  ls -l
```

You can manage your shell history with the `history` command and the `HISTFILE` environment variable:

- Read history from a file with `-r <histfile>`:
```bash
$ history -r <path_to_history_file>
$ history
1  history -r <path_to_history_file>
2  echo hello
3  echo world
4  history
```

- Save your current session's history with `-w <histfile>`:
```bash
$ echo hello
hello
$ echo world
world
$ history -w <path_to_history_file>
```

- Append history to a file `-a <histfile>`:
```bash
$ echo new_command
new_command
$ history -a <path_to_history_file>
```

- Automatically read from and append to a history file on shell startup and exit with `HISTFILE` environment variable:
```bash
HISTFILE=<path_to_history_file> ./your_program.sh
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