package reader

import (
	"bufio"
	"slices"
	"strings"
)

type Reader struct {
	Reader *bufio.Reader
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

	input = strings.TrimSpace(input)
	tokens := tokenize(input)

	if len(tokens) == 0 {
		return Command{}, nil
	}

	return Command{
		Name: strings.ToLower(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func tokenize(input string) []string {
	var (
		tokens      []string
		token       string
		openQuote   rune
		escapeMode  bool
		wasInQuotes bool
		redirectors []string = []string{">", "1>", ">>", "1>>", "2>", "2>>"}
	)

	for _, ch := range input {
		switch {
		case escapeMode:
			if openQuote != 0 && openQuote != '\'' && (openQuote != '"' || (ch != '\\' && ch != '$' && ch != '"')) {
				token += "\\"
			}
			token += string(ch)
			escapeMode = false
		case ch == '\\' && openQuote != '\'':
			escapeMode = true
		case ch == openQuote:
			openQuote = 0
		case ch == '\'' || ch == '"':
			if openQuote == 0 {
				openQuote = ch
				wasInQuotes = true
			} else {
				token += string(ch)
			}
		case ch == ' ' && openQuote == 0:
			// Detect output redirect
			if !wasInQuotes && slices.Contains(redirectors, token) {
				token = "'" + token + "'"
			}

			if len(token) > 0 {
				tokens = append(tokens, token)
				token = ""
				wasInQuotes = false
			}
		default:
			token += string(ch)
		}
	}

	if len(token) > 0 {
		tokens = append(tokens, token)
	}

	return tokens
}
