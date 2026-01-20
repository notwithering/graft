package syntax

import "regexp"

type Syntax struct {
	OpenClose *regexp.Regexp
	Close     *regexp.Regexp
	Parse     func(string) ([]string, error)
}
