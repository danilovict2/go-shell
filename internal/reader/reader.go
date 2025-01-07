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

	s := strings.Trim(command, "\r\n")
	tokens := make([]string, 0)

	for {
		start := strings.Index(s, "'")
		if start == -1 {
			tokens = append(tokens, strings.Fields(s)...)
			break
		}

		tokens = append(tokens, strings.Fields(s[:start])...)
		s = s[start+1:]
		end := strings.Index(s, "'")
		token := s[:end]
		tokens = append(tokens, token)
		s = s[end+1:]
	}

	return Command{
		Name: strings.ToLower(tokens[0]),
		Args: tokens[1:],
	}, nil
}
