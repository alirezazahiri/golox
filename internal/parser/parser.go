package parser

import (
	"golox/internal/chunk"
	"golox/internal/scanner"
	"golox/internal/vm"
)

type Parser struct {
	Current   scanner.Token
	Previous  scanner.Token
	HadError  bool
	PanicMode bool
	scanner   *scanner.Scanner
	chunk     *chunk.Chunk
	vm        *vm.VM
	rules     map[scanner.TokenType]ParserRule
}

func New(s *scanner.Scanner, c *chunk.Chunk, v *vm.VM) *Parser {
	p := &Parser{
		scanner: s,
		chunk:   c,
		vm:      v,
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
