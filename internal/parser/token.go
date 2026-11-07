package parser

import "golox/internal/scanner"

func identifiersEqual(a *scanner.Token, b *scanner.Token) bool {
	return a.Lexeme == b.Lexeme
}
