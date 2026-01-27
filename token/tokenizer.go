package token

import (
	"fmt"
)

func Tokenize(data string, syntax *Syntax) ([]*Token, error) {
	var tokens []*Token

	cursor := 0

	for _, match := range syntax.OpenClose.FindAllStringIndex(data, -1) {
		fullStart, fullEnd := match[0], match[1]
		// subStart, subEnd := match[2], match[3]

		rawMatch := data[fullStart:fullEnd]
		// rawCommand := data[subStart:subEnd]

		if fullStart > cursor {
			tokens = append(tokens, &Token{
				Kind: TokenText,
				Text: data[cursor:fullStart],
			})
		}

		if syntax.Close.MatchString(rawMatch) {
			tokens = append(tokens, &Token{
				Kind: TokenClose,
				Text: rawMatch,
			})

			cursor = fullEnd
			continue
		}

		command, data, err := syntax.Parse(rawMatch)
		if err != nil {
			return nil, fmt.Errorf("parsing args: %w", err)
		}

		if command == "" {
			cursor = fullEnd
			continue
		}

		tokens = append(tokens, &Token{
			Kind:    TokenOpen,
			Command: command,
			Data:    data,
			Text:    rawMatch,
		})

		cursor = fullEnd
	}

	if cursor < len(data) {
		tokens = append(tokens, &Token{
			Kind: TokenText,
			Text: data[cursor:],
		})
	}

	return tokens, nil
}
