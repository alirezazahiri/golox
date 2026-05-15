package scanner

import (
	"golox/pkg/errors"
)

func (s *Scanner) string() Token {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.Line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return s.errorToken(errors.UnterminatedStringErrorMessage)
	}

	s.advance()

	return s.makeToken(TOKEN_STRING)
}

func (s *Scanner) number() Token {
	for isDigit(s.peek()) {
		s.advance()
	}

	// fractional part
	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance() // '.'
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	return s.makeToken(TOKEN_NUMBER)
}