package reader

import (
	"fmt"
	"io"
	"os"
	"slices"
)

type Command struct {
	Name string
	Args []string
}

func (c Command) String() string {
	return c.Name
}

func (c *Command) GetOutputWriters() (stdout io.Writer, stderr io.Writer, err error) {
    stdout = os.Stdout
	stderr = os.Stderr
	
	for i := 0; i < len(c.Args)-1; i++ {
		switch c.Args[i] {
		case ">", "1>":
			stdout, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_CREATE, 0644)
            if err != nil {
                return nil, nil, fmt.Errorf("error opening file: %w", err)
            }

			c.Args = slices.Delete(c.Args, i, i+2)
		case "2>":
			stderr, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_CREATE, 0644)
            if err != nil {
                return nil, nil, fmt.Errorf("error opening file: %w", err)
            }

			c.Args = slices.Delete(c.Args, i, i+2)
		}
    }

    return stdout, stderr, nil
}