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
	input, err := r.Reader.ReadString('\n')
	if err != nil {
		return Command{}, err
	}

	input = strings.Trim(input, "\r\n")
	var (
		tokens    []string = make([]string, 0)
		token     string
		openQuote rune
	)

	for _, ch := range input {
		switch {
		case ch == openQuote:
			openQuote = 0
		case ch == '\'' || ch == '"':
			if openQuote == 0 {
				openQuote = ch
			} else {
				token += string(ch)
			}
		case ch == ' ' && openQuote == 0:
			if token != "" {
				tokens = append(tokens, token)
				token = ""
			}
		default:
			token += string(ch)
		}
	}

	if token != "" {
		tokens = append(tokens, token)
	}

	return Command{
		Name: strings.ToLower(tokens[0]),
		Args: tokens[1:],
	}, nil
}
