package reader

import (
	"bufio"
	"strings"
)

type Reader struct {
	Reader *bufio.Reader
}

type Command struct {
	Name string
	Args []string
}

func (c Command) String() string {
	return c.Name
}

func New(r *bufio.Reader) Reader {
	return Reader{
		Reader: r,
	}
}

func (r Reader) Read() (Command, error) {
	command, err := r.Reader.ReadString('\n')
	if err != nil {
		return Command{}, err
	}

	command = command[:len(command)-1]
	split := strings.Split(command, " ")

	ret := Command{
		Name: split[0],
		Args: make([]string, 0),
	}

	for i := 1; i < len(split); i++ {
		ret.Args = append(ret.Args, split[i])
	}

	return ret, nil
}
