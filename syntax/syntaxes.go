package syntax

import (
	"regexp"
	"strings"

	"github.com/mattn/go-shellwords"
)

var (
	SGMLCommentSyntax = &Syntax{
		OpenClose: regexp.MustCompile(`(?s)<!--(.+?)-->`),
		Close:     regexp.MustCompile(`(?s)<!--\s*end\s*-->`),
		Parse:     shellwords.Parse,
	}
	SGMLTagSyntax = &Syntax{
		OpenClose: regexp.MustCompile(`</?g-(.+?)\s*/?>`),
		Close:     regexp.MustCompile(`</g-.+?>`),
		Parse:     shellwords.Parse,
	}
	// replaces backslashes with forward slashes, so prettier doesnt shit the bed when seeing a path
	PrettierHTMLTagSyntax = &Syntax{
		OpenClose: regexp.MustCompile(`</?g-(.+?)\s*/?>`),
		Close:     regexp.MustCompile(`</g-.+?>`),
		Parse: func(s string) ([]string, error) {
			return shellwords.Parse(strings.ReplaceAll(s, "\\", "/"))
		},
	}
	DoubleBraceSyntax = &Syntax{
		OpenClose: regexp.MustCompile(`\{\{(.+?)\}\}`),
		Close:     regexp.MustCompile(`\{\{end\}\}`),
		Parse:     shellwords.Parse,
	}
	DunderSyntax = &Syntax{
		OpenClose: regexp.MustCompile(`__(.+?)__`),
		Close:     regexp.MustCompile(`__end__`),
		Parse:     shellwords.Parse,
	}
)
