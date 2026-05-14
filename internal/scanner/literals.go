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

	// value := s.source[s.Start+1 : s.Current-1]
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

	// val, _ := strconv.ParseFloat(s.source[s.Start:s.Current], 64)
	return s.makeToken(TOKEN_NUMBER)
}

func (s *Scanner) addLiteral(tokenType TokenType, literal any) {
	// text := s.source[s.Start:s.Current]
	// s.tokens = append(s.tokens, NewToken(
	// 	tokenType,
	// 	text,
	// 	literal,
	// 	s.line,
	// ))
}
