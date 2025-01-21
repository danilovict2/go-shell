package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/executable"
	"golang.org/x/term"
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
	input, err := p.readInput()
	if err != nil {
		return command.Command{}, err
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

func (p Parser) readInput() (string, error) {
	var (
		input     string
		doubletab bool
	)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}

	defer func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
		fmt.Fprintf(os.Stdout, "\n")
	}()

Loop:
	for {
		b, err := p.Reader.ReadByte()
		if err != nil {
			break Loop
		}

		switch b {
		case '\r', '\n':
			break Loop
		case '\t':
			suffixes := autocomplete(input)
			switch {
			case len(suffixes) == 1:
				doubletab = false
				suffix := suffixes[0] + " "
				input += suffix
				fmt.Fprint(os.Stdout, suffix)
			case len(suffixes) > 1:
				if doubletab {
					fmt.Fprint(os.Stdout, "\r\n")
					for _, suffix := range suffixes {
						fmt.Fprint(os.Stdout, input, suffix, "  ")
					}

					fmt.Fprint(os.Stdout, "\r\n$ ", input)
				} else {
					fmt.Fprint(os.Stdout, "\a")
				}

				doubletab = !doubletab
			default:
				doubletab = false
				fmt.Fprint(os.Stdout, "\a")
			}
		case '\x7F':
			if len(input) > 0 {
				input = input[:len(input)-1]
				fmt.Fprint(os.Stdout, "\b \b")
			}
		case '\x03':
			fmt.Fprintf(os.Stdout, "^C")
			return "", fmt.Errorf("^C")
		default:
			if b >= 32 {
				input += string(b)
				fmt.Fprint(os.Stdout, string(b))
			}
		}
	}

	return input, err
}

func autocomplete(prefix string) (suffixes []string) {
	suffixes = make([]string, 0)
	if len(prefix) == 0 {
		return suffixes
	}

	for _, command := range command.Builtins {
		if strings.HasPrefix(command, prefix) {
			suffixes = append(suffixes, command[len(prefix):])
		}
	}

	for _, command := range executable.Executables {
		command = filepath.Base(command)
		var suffix string

		if len(command) >= len(prefix) {
			suffix = command[len(prefix):]
		}

		if strings.HasPrefix(command, prefix) && !slices.Contains(suffixes, suffix) {
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
