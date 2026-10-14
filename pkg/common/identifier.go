package common

import "golox/internal/scanner"

func IdentifiersEqual(a *scanner.Token, b *scanner.Token) bool {
	return a.Lexeme == b.Lexeme
}
