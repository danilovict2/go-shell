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

func (p Parser) ParseInput() ([]command.Command, error) {
	line, err := p.readInput()
	if err != nil {
		return nil, err
	}

	commands := strings.Split(line, "|")
	ret := make([]command.Command, 0)

	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		tokens := tokenize(cmd)
		if len(tokens) == 0 {
			continue
		}

		ret = append(ret, command.New(strings.ToLower(tokens[0]), tokens[1:]))
	}

	return ret, nil
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
			input, doubletab = handleTab(input, doubletab)
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

func handleTab(input string, doubletab bool) (string, bool) {
	suffixes := autocomplete(input)
	switch len(suffixes) {
	case 1:
		input = appendSuffix(input, suffixes[0]+" ")
		doubletab = false
	case 0:
		doubletab = false
		fmt.Fprint(os.Stdout, "\a")
	default:
		if allHaveSamePrefix(suffixes) {
			input = appendSuffix(input, suffixes[0])
		} else {
			if doubletab {
				fmt.Fprint(os.Stdout, "\r\n")
				for _, suffix := range suffixes {
					fmt.Fprint(os.Stdout, input, suffix, "  ")
				}
				fmt.Fprint(os.Stdout, "\r\n$ ", input)
			} else {
				fmt.Fprint(os.Stdout, "\a")
			}
		}

		doubletab = !doubletab
	}
	return input, doubletab
}

func appendSuffix(input, suffix string) string {
	input += suffix
	fmt.Fprint(os.Stdout, suffix)
	return input
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

	executables := executable.FindExecutables()
	for _, command := range executables {
		command = filepath.Base(command)
		var suffix string

		if len(command) >= len(prefix) {
			suffix = command[len(prefix):]
		}

		if strings.HasPrefix(command, prefix) && !slices.Contains(suffixes, suffix) {
			suffixes = append(suffixes, command[len(prefix):])
		}
	}

	slices.Sort(suffixes)
	return suffixes
}

func allHaveSamePrefix(suffixes []string) bool {
	for _, suffix := range suffixes {
		if !strings.HasPrefix(suffix, suffixes[0]) {
			return false
		}
	}

	return true
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
