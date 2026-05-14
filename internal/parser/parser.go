package parser

import (
	"golox/internal/chunk"
	"golox/internal/scanner"
)

type Parser struct {
	Current   scanner.Token
	Previous  scanner.Token
	HadError  bool
	PanicMode bool
	scanner   *scanner.Scanner
	chunk     *chunk.Chunk
	rules     map[scanner.TokenType]ParserRule
}

func New(s *scanner.Scanner, c *chunk.Chunk) *Parser {
	p := &Parser{
		scanner: s,
		chunk:   c,
	}
	p.InitParserRules()
	return p
}

func (p *Parser) Advance() {
	p.Previous = p.Current

	for {
		p.Current = p.scanner.ScanToken()
		if p.Current.Type != scanner.TOKEN_ERROR {
			break
		}

		p.errorAtCurrent(p.Current.Lexeme)
	}
}

func (p *Parser) Consume(t scanner.TokenType, msg string) {
	if p.Current.Type == t {
		p.Advance()
		return
	}

	p.errorAtCurrent(msg)
}
