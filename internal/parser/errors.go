package parser

import (
	"fmt"
	"golox/internal/scanner"
)

func (p *Parser) errorAtCurrent(msg string) {
	p.errorAt(p.Current, msg)
}

func (p *Parser) error(msg string) {
	p.errorAt(p.Previous, msg)
}

func (p *Parser) errorAt(token scanner.Token, msg string) {
	if p.PanicMode {
		return
	}
	p.PanicMode = true
	fmt.Printf("[line %d] Error", token.Line)

	switch token.Type {
	case scanner.TOKEN_EOF:
		fmt.Print(" at the end")
	case scanner.TOKEN_ERROR: // noop!
		break
	default:
		fmt.Printf(" at '%s'", token.Lexeme)
	}

	fmt.Printf(": %s\n", msg)
	p.HadError = true
}
