package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/executable"
)

type Command struct {
	Name  string
	Args  []string
	Input string
}

func (c Command) String() string {
	return c.Name
}

func (c *Command) GetOutputWriters() (stdout io.Writer, stderr io.Writer, err error) {
	stdout = os.Stdout
	stderr = os.Stderr

	for i := 0; i < len(c.Args)-1; i++ {
		switch c.Args[i] {
		case "'>'", "'1>'":
			stdout, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return nil, nil, fmt.Errorf("error opening file: %w", err)
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		case "'>>'", "'1>>'":
			stdout, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				return nil, nil, fmt.Errorf("error opening file: %w", err)
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		case "'2>'":
			stderr, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return nil, nil, fmt.Errorf("error opening file: %w", err)
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		case "'2>>'":
			stderr, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
			if err != nil {
				return nil, nil, fmt.Errorf("error opening file: %w", err)
			}

			c.Args = slices.Delete(c.Args, i, i+2)
		}
	}

	return stdout, stderr, nil
}

func Pipeline(commands []Command) (Command, error) {
	for i := range len(commands) - 1 {
		out, err := commands[i].GetOutput()
		if err != nil {
			return Command{}, err
		}

		commands[i+1].Input = out
	}

	return commands[len(commands)-1], nil
}

func (c *Command) GetOutput() (string, error) {
	handler, isBuiltin := BuiltinHandlers[c.Name]
	if isBuiltin {
		return handler(c.Args), nil
	}

	out, err := c.getNonBuiltinOutput()
	if err != nil {
		return "", err
	}

	return strings.TrimRight(string(out), "\n"), nil
}

func (c *Command) getNonBuiltinOutput() ([]byte, error) {
	if executable.GetExecutableFilePath(c.Name) == "" {
		return nil, fmt.Errorf("%s: command not found", c.Name)
	}

	comm := exec.Command(c.Name, c.Args...)
	stdin, err := comm.StdinPipe()
	if err != nil {
		return nil, err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, c.Input)
	}()

	var stderrBuf bytes.Buffer
	comm.Stderr = &stderrBuf

	out, err := comm.Output()
	if err != nil {
		msg := strings.TrimSpace(stderrBuf.String())
		return out, fmt.Errorf("%s", msg)
	}

	return out, nil
}
