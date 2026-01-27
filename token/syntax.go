package token

import "regexp"

type Syntax struct {
	OpenClose *regexp.Regexp
	Close     *regexp.Regexp
	Parse     func(string) (command string, data any, err error)
}
