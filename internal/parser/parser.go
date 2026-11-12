package parser

import (
	"golox/internal/scanner"
	"golox/pkg/bytecode"
)

type Parser struct {
	Current   scanner.Token
	Previous  scanner.Token
	HadError  bool
	PanicMode bool
	scanner   *scanner.Scanner
	chunk     *bytecode.Chunk
	allocator Allocator
	rules     map[scanner.TokenType]ParseRule
	Compiler  *Compiler
}

func New(s *scanner.Scanner, allocator Allocator) *Parser {
	compiler := NewCompilerWithFunction(TypeScript, allocator.AllocateFunction())
	p := &Parser{
		scanner:   s,
		chunk:     compiler.Function.Chunk,
		allocator: allocator,
		Compiler:  compiler,
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

func (p *Parser) Match(t scanner.TokenType) bool {
	if !p.Check(t) {
		return false
	}

	p.Advance()
	return true
}

func (p *Parser) Check(t scanner.TokenType) bool {
	return p.Current.Type == t
}
