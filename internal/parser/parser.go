package parser

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/internal/command"
	"github.com/codecrafters-io/shell-starter-go/internal/history"
	"github.com/codecrafters-io/shell-starter-go/internal/parameter"
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

func (p Parser) ParseInput() (ret []command.Command, err error) {
	line, err := p.readInput()
	if err != nil {
		return nil, err
	}

	history.Commands = append(history.Commands, line)

	background := false
	if line[len(line)-1] == '&' {
		background = true
		line = line[:len(line)-1]
	}

	commands := strings.Split(line, "|")
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		tokens := tokenize(cmd)
		if len(tokens) == 0 {
			continue
		}

		ret = append(ret, command.New(strings.ToLower(tokens[0]), parseArgs(tokens[1:]), background))
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

	historyPos := len(history.Commands)
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
		case 0x1b:
			bracket, err := p.Reader.ReadByte()
			if err != nil || bracket != '[' {
				break Loop
			}

			arrowCode, err := p.Reader.ReadByte()
			if err != nil {
				break Loop
			}

			switch arrowCode {
			case 'A':
				if len(history.Commands) == 0 {
					bell()
					continue
				}

				if historyPos > 0 {
					historyPos--
					if historyPos < len(history.Commands) {
						input = history.Commands[historyPos]
						clearLine()
						fmt.Fprintf(os.Stdout, "%s", input)
					} else {
						bell()
						historyPos--
					}
				} else {
					bell()
				}
			case 'B':
				if len(history.Commands) == 0 || historyPos < 0 {
					bell()
					continue
				}

				historyPos++
				if historyPos < len(history.Commands) {
					input = history.Commands[historyPos]
					clearLine()
					fmt.Fprintf(os.Stdout, "%s", input)
				} else {
					bell()
					historyPos--
				}
			default:
				break
			}

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
	completions := autocomplete(input)
	prefix := commonPrefix(completions)

	switch {
	case len(completions) == 1:
		input += completions[0].String()
		fmt.Fprint(os.Stdout, completions[0])
		doubletab = false
	case len(completions) == 0:
		doubletab = false
		bell()
	case prefix != "":
		input += prefix
		fmt.Fprint(os.Stdout, prefix)
		doubletab = false
	default:
		if doubletab {
			fmt.Fprint(os.Stdout, "\r\n")
			for _, completion := range completions {
				fmt.Fprint(os.Stdout, completion.Prefix, completion)
				if completion.Trailing != " " {
					fmt.Fprint(os.Stdout, " ")
				}
			}
			fmt.Fprint(os.Stdout, "\r\n$ ", input)
		} else {
			bell()
		}

		doubletab = !doubletab
	}
	return input, doubletab
}

func clearLine() {
	fmt.Fprintf(os.Stdout, "\r\x1b[D")
	fmt.Fprint(os.Stdout, "\r\033[K")
	fmt.Fprint(os.Stdout, "\x1b[D")
	fmt.Fprint(os.Stdout, "$ ")
}

func bell() {
	fmt.Fprint(os.Stdout, "\a")
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

func parseArgs(args []string) []string {
	parsedArgs := make([]string, 0)

	for i, arg := range args {
		for j := range arg {
			if arg[j] == '$' {
				varName := arg[j+1:]
				k := len(arg)
				if j+1 < len(arg) && arg[j+1] == '{' {
					if idx := strings.Index(varName, "}"); idx != -1 {
						idx += j
						varName = arg[j+2 : idx+1]
						k = idx + 2
					}
				}

				val, _ := parameter.Get(varName)
				args[i] = arg[:j] + val + arg[k:]
			}
		}

		if args[i] != "" {
			parsedArgs = append(parsedArgs, args[i])
		}
	}

	return parsedArgs
}
