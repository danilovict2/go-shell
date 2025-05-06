package command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"

	"github.com/codecrafters-io/shell-starter-go/internal/executable"
)


type Command struct {
	Name       string
	Args       []string
	Stdout     io.WriteCloser
	Stderr     io.WriteCloser
	OutputChan chan message
}

type message struct {
	data  []byte
	isError bool
}

func New(name string, args []string) *Command {
	c := &Command{
		Name:       name,
		Args:       args,
		OutputChan: make(chan message),
	}

	c.setOutputWriters()
	return c
}

func (c Command) String() string {
	return c.Name
}

func (c *Command) setOutputWriters() {
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	var err error

	for i := range len(c.Args) - 1 {
		switch c.Args[i] {
		case "'>'", "'1>'":
			c.Stdout, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening file: %s\n", err.Error())
				return
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		case "'>>'", "'1>>'":
			c.Stdout, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening file: %s\n", err.Error())
				return
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		case "'2>'":
			c.Stderr, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening file: %s\n", err.Error())
				return
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		case "'2>>'":
			c.Stderr, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error opening file: %s\n", err.Error())
				return
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		}
	}
}

func Pipeline(commands []Command) {
	go commands[0].execute(*New("", nil))
	isLastNewLine := true

	for out := range commands[0].OutputChan {
		data := string(out.data)
		isLastNewLine = data[len(data)-1] == '\n'

		if out.isError {
			fmt.Fprint(commands[0].Stderr, data)
		} else {
			fmt.Fprint(commands[0].Stdout, data)
		}
	}

	if !isLastNewLine {
		fmt.Fprintln(commands[0].Stdout)
	}
}

func (c *Command) execute(prev Command) {
	handler, isBuiltin := BuiltinHandlers[c.Name]
	if isBuiltin {
		if output := handler(c.Args); output != "" {
			c.OutputChan <- message{data: []byte(output), isError: false}
		}

		close(c.OutputChan)
		return
	}

	go c.executeNonBuiltin(prev)
}

func (c *Command) executeNonBuiltin(prev Command) {
	defer close(c.OutputChan)
	if executable.GetExecutableFilePath(c.Name) == "" {
		c.OutputChan <- message{data: fmt.Appendf(nil, "%s: command not found\n", c.Name), isError: true}
		return
	}

	comm := exec.Command(c.Name, c.Args...)
	stdin, err := comm.StdinPipe()
	if err != nil {
		return
	}

	stdout, err := comm.StdoutPipe()
	if err != nil {
		return
	}

	stderr, err := comm.StderrPipe()
	if err != nil {
		return
	}

	if err := comm.Start(); err != nil {
		return
	}

	go func() {
		for input := range prev.OutputChan {
			if !input.isError {
				stdin.Write(input.data)
			}
		}
	}()

	go func() {
		data := make([]byte, 1024)
		for {
			n, err := stdout.Read(data)
			if err != nil {
				break
			}

			c.OutputChan <- message{data: data[:n], isError: false}
		}
	}()

	slurp, _ := io.ReadAll(stderr)
	if len(slurp) > 0 {
		c.OutputChan <- message{data: slurp, isError: true}
	}

	comm.Wait()
}
