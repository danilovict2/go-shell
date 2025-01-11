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

func (c *Command) GetOutputWriter() (io.Writer, error) {
	stdout := os.Stdout
    var err error

    for i := 0; i < len(c.Args)-1; i++ {
        if c.Args[i] == ">" || c.Args[i] == "1>" {
            stdout, err = os.OpenFile(c.Args[i+1], os.O_WRONLY|os.O_CREATE, 0644)
            if err != nil {
                return nil, fmt.Errorf("error opening file: %w", err)
            }

			c.Args = slices.Delete(c.Args, i, i+2)
			break
        }
    }

    return stdout, nil
}