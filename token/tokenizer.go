package token

import (
	"fmt"

	"github.com/notwithering/graft/syntax"
)

func Tokenize(data string, syntax *syntax.Syntax) ([]*Token, error) {
	var tokens []*Token

	cursor := 0

	for _, match := range syntax.OpenClose.FindAllStringSubmatchIndex(data, -1) {
		fullStart, fullEnd := match[0], match[1]
		subStart, subEnd := match[2], match[3]

		rawMatch := data[fullStart:fullEnd]
		rawCommand := data[subStart:subEnd]

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

		args, err := syntax.Parse(rawCommand)
		if err != nil {
			return nil, fmt.Errorf("parsing args: %w", err)
		}

		if len(args) == 0 {
			cursor = fullEnd
			continue
		}

		tokens = append(tokens, &Token{
			Kind: TokenOpen,
			Args: args,
			Text: rawMatch,
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
