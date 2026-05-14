package scanner

import (
	"golox/pkg/errors"
)

func (s *Scanner) ScanToken() Token {
	s.skipWhitespace()
	s.Start = s.Current

	if s.isAtEnd() {
		return s.makeToken(TOKEN_EOF)
	}

	c := s.advance()

	switch c {
	case '(':
		return s.makeToken(TOKEN_LEFT_PAREN)
	case ')':
		return s.makeToken(TOKEN_RIGHT_PAREN)
	case '{':
		return s.makeToken(TOKEN_LEFT_BRACE)
	case '}':
		return s.makeToken(TOKEN_RIGHT_BRACE)
	case ';':
		return s.makeToken(TOKEN_SEMICOLON)
	case ',':
		return s.makeToken(TOKEN_COMMA)
	case '.':
		return s.makeToken(TOKEN_DOT)
	case '-':
		return s.makeToken(TOKEN_MINUS)
	case '+':
		return s.makeToken(TOKEN_PLUS)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
			return s.ScanToken()
		}
		return s.makeToken(TOKEN_SLASH)
	case '*':
		return s.makeToken(TOKEN_STAR)
	case '!':
		if s.match('=') {
			return s.makeToken(TOKEN_BANG_EQUAL)
		}
		return s.makeToken(TOKEN_BANG)
	case '=':
		if s.match('=') {
			return s.makeToken(TOKEN_EQUAL_EQUAL)
		}
		return s.makeToken(TOKEN_EQUAL)
	case '<':
		if s.match('=') {
			return s.makeToken(TOKEN_LESS_EQUAL)
		}
		return s.makeToken(TOKEN_LESS)
	case '>':
		if s.match('=') {
			return s.makeToken(TOKEN_GREATER_EQUAL)
		}
		return s.makeToken(TOKEN_GREATER)
	case '\n':
		s.Line++
		return s.ScanToken()
	case '"':
		return s.string()
	default:
		if isDigit(c) {
			return s.number()
		} else if isAlpha(c) {
			return s.identifier()
		}
		break
	}

	return s.errorToken(errors.UnexpectedCharacterErrorMessage)
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.Current]
}

func (s *Scanner) peekNext() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.Current+1]
}

func (s *Scanner) advance() byte {
	s.Current++
	return s.source[s.Current-1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.Current] != expected {
		return false
	}

	s.Current++
	return true
}

func (s *Scanner) isAtEnd() bool {
	return s.Current >= s.sourceLength
}

func (s *Scanner) makeToken(t TokenType) Token {
	return NewToken(t, s.Start, s.Current-s.Start, s.Line, s.source[s.Start:s.Current])
}

func (s *Scanner) errorToken(msg string) Token {
	return NewToken(TOKEN_ERROR, 0, len(msg), s.Line, msg)
}

func (s *Scanner) skipWhitespace() {
	for {
		if s.isAtEnd() {
			return
		}

		c := s.peek()

		switch c {
		case ' ', '\r', '\t':
			s.advance()

		case '\n':
			s.Line++
			s.advance()

		case '/':
			if s.peekNext() == '/' {
				for s.peek() != '\n' && !s.isAtEnd() {
					s.advance()
				}
			} else {
				return
			}

		default:
			return
		}
	}
}

func (s *Scanner) identifier() Token {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	return s.makeToken(s.identifierType())
}

func (s *Scanner) identifierType() TokenType {
	text := s.source[s.Start:s.Current]
	if tok, ok := keywords[text]; ok {
		return tok
	}
	return TOKEN_IDENTIFIER
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
