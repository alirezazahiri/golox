package scanner

type Token struct {
	Type   TokenType
	Start  int
	Length int
	Line   int
	Lexeme string
}

func NewToken(t TokenType, start, length, line int, lexeme string) Token {
	return Token{
		Type:   t,
		Start:  start,
		Length: length,
		Line:   line,
		Lexeme: lexeme,
	}
}
