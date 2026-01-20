package syntax

import (
	"regexp"
)

var (
	SGMLCommentSyntax = &Syntax{
		OpenClose: regexp.MustCompile(`(?s)<!--(.+?)-->`),
		Close:     regexp.MustCompile(`(?s)<!--\s*end\s*-->`),
	}
	SGMLTagSyntax = &Syntax{
		OpenClose: regexp.MustCompile(`</?g-(.+?)\s*/?>`),
		Close:     regexp.MustCompile(`</g-.+?>`),
	}
	DoubleBraceSyntax = &Syntax{
		OpenClose: regexp.MustCompile(`\{\{(.+?)\}\}`),
		Close:     regexp.MustCompile(`\{\{end\}\}`),
	}
	DunderSyntax = &Syntax{
		OpenClose: regexp.MustCompile(`__(.+?)__`),
		Close:     regexp.MustCompile(`__end__`),
	}
)
