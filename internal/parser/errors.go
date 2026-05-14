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

	if token.Type == scanner.TOKEN_EOF {
		fmt.Print(" at end")
	} else if token.Type == scanner.TOKEN_ERROR {
		// noop!
	} else {
		fmt.Printf(" at '%s'", token.Lexeme)
	}

	fmt.Printf(": %s\n", msg)
	p.HadError = true
}
