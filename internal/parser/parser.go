package parser

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
)

type Parser struct {
	Reader *bufio.Reader
}

func New(r *bufio.Reader) Parser {
	return Parser{
		Reader: r,
	}
}

func (p Parser) ParseInput() (command.Command, error) {
	input := ""
Loop:
	for {
		b, err := p.Reader.ReadByte()
		if err != nil {
			return command.Command{}, err
		}

		switch b {
		case '\r':
			fmt.Fprint(os.Stdout, "\r\n")
			break Loop
		case '\t':
			suffixes := autocomplete(input)
			if len(suffixes) > 0 {
				suffix := suffixes[0] + " "
				input += suffix
				fmt.Fprint(os.Stdout, suffix)
			}
		case '\x7F':
			if len(input) > 0 {
				input = input[:len(input)-1]
				fmt.Fprint(os.Stdout, "\b \b")
			}
		case '\x03':
			return command.Command{}, fmt.Errorf("Ctrl+C")
		default:
			input += string(b)
			fmt.Fprint(os.Stdout, string(b))
		}
	}

	input = strings.TrimSpace(input)
	tokens := tokenize(input)

	if len(tokens) == 0 {
		return command.Command{}, nil
	}

	return command.Command{
		Name: strings.ToLower(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func autocomplete(prefix string) (suffixes []string) {
	suffixes = make([]string, 0)

	for _, command := range command.Builtins {
		if strings.HasPrefix(command, prefix) {
			suffixes = append(suffixes, command[len(prefix):])
		}
	}

	return suffixes
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
