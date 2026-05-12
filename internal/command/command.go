package command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"sync"

	"github.com/codecrafters-io/shell-starter-go/internal/executable"
)

type Command struct {
	Name string
	Args []string
}

func New(name string, args []string) Command {
	return Command{
		Name: name,
		Args: args,
	}
}

func (c Command) String() string {
	return c.Name
}

func (c *Command) getOutputWriters() (stdout, stderr io.WriteCloser, err error) {
	stdout = os.Stdout
	stderr = os.Stderr

	for i := range len(c.Args) - 1 {
		switch c.Args[i] {
		case "'>'", "'1>'":
			stdout, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return nil, nil, fmt.Errorf("error opening file: %v", err)
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		case "'>>'", "'1>>'":
			stdout, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				return nil, nil, fmt.Errorf("error opening file: %v", err)
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		case "'2>'":
			stderr, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return nil, nil, fmt.Errorf("error opening file: %v", err)
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		case "'2>>'":
			stderr, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				return nil, nil, fmt.Errorf("error opening file: %v", err)
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		}
	}

	return stdout, stderr, nil
}

func Pipeline(commands []Command) {
	if len(commands) == 0 {
		return
	}

	stderrs := make([]io.WriteCloser, 0)
	var (
		stdout, stderr io.WriteCloser
		err            error
	)

	for i, command := range commands {
		stdout, stderr, err = command.getOutputWriters()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		stderrs = append(stderrs, stderr)
		commands[i] = command
	}

	var (
		reader io.ReadCloser = os.Stdin
		wg                   = &sync.WaitGroup{}
		writer io.WriteCloser
		pr     io.ReadCloser
		pw     io.WriteCloser
	)

	for i, cmd := range commands {
		if i == len(commands)-1 {
			writer = stdout
		} else {
			pr, pw = io.Pipe()
			writer = pw
		}

		wg.Add(1)
		go cmd.execute(reader, writer, stderrs[i], wg)
		reader = pr
	}

	wg.Wait()
}

func (c *Command) execute(stdin io.ReadCloser, stdout, stderr io.WriteCloser, wg *sync.WaitGroup) {
	handler, isBuiltin := BuiltinHandlers[c.Name]
	if isBuiltin {
		if output, err := handler(c.Args); err != nil {
			fmt.Fprintln(stderr, err)
		} else if output != "" {
			fmt.Fprintln(stdout, output)
		}
	} else {
		c.executeNonBuiltin(stdin, stdout, stderr)
	}

	for _, f := range []io.Closer{stdin, stdout, stderr} {
		if f != os.Stdin && f != os.Stdout && f != os.Stderr {
			f.Close()
		}
	}

	if wg != nil {
		wg.Done()
	}
}

func (c *Command) executeNonBuiltin(stdin io.Reader, stdout, stderr io.Writer) {
	if executable.GetExecutableFilePath(c.Name) == "" {
		fmt.Fprintf(stderr, "%s: command not found\n", c.Name)
		return
	}

	comm := exec.Command(c.Name, c.Args...)
	comm.Stdin = stdin
	comm.Stdout = stdout
	comm.Stderr = stderr

	comm.Run()
}
