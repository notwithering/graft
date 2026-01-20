package token

import "fmt"

type Token struct {
	Kind TokenKind
	Args []string
	Text string
}

type TokenKind uint8

const (
	TokenText TokenKind = iota
	TokenOpen
	TokenClose
)

func (t Token) String() string {
	switch t.Kind {
	case TokenText:
		return fmt.Sprintf("TokenText(%q)", t.Text)
	case TokenOpen:
		return fmt.Sprintf("TokenOpen(%q)", t.Text)
	case TokenClose:
		return fmt.Sprintf("TokenClose(%q)", t.Text)
	default:
		return fmt.Sprintf("Token(%d)", t.Kind)
	}
}
