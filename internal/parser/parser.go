package parser

import (
	"golox/internal/scanner"
	"golox/internal/vm"
	"golox/pkg/bytecode"
)

type Parser struct {
	Current   scanner.Token
	Previous  scanner.Token
	HadError  bool
	PanicMode bool
	scanner   *scanner.Scanner
	chunk     *bytecode.Chunk
	vm        *vm.VM
	rules     map[scanner.TokenType]ParseRule
	Compiler  *Compiler
}

func New(s *scanner.Scanner, c *bytecode.Chunk, v *vm.VM) *Parser {
	p := &Parser{
		scanner:  s,
		chunk:    c,
		vm:       v,
		Compiler: NewCompiler(TypeScript),
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
