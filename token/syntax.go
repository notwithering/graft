package token

import "regexp"

type Syntax struct {
	// OpenClose should match both open and close tokens.
	OpenClose *regexp.Regexp

	// Close should match only the close token.
	Close *regexp.Regexp

	// Parse should parse the token and return the command name and data.
	Parse func(string) (command string, data any, err error)
}
