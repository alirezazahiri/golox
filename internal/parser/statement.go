package parser

import (
	"golox/internal/scanner"
	"golox/pkg/common"
	"golox/pkg/errors"
)

func (p *Parser) Declaration() {

	if p.Match(scanner.TOKEN_VAR) {
		p.varDeclaration()
	} else {
		p.statement()
	}

	if p.PanicMode {
		p.synchronize()
	}
}

func (p *Parser) varDeclaration() {
	global := p.parseVariable(errors.ExpectIdentifierAfterVar)
	
	if p.Match(scanner.TOKEN_EQUAL) {
		p.expression()
	} else {
		p.EmitByte(byte(common.OpNil))
	}

	p.Consume(scanner.TOKEN_SEMICOLON, errors.ExpectSemicolonAfterVarDeclaration)
	p.defineVariable(global)
}

func (p *Parser) statement() {
	if p.Match(scanner.TOKEN_PRINT) {
		p.printStatement()
	} else {
		p.expressionStatement()
	}
}

func (p *Parser) printStatement() {
	p.expression()
	p.Consume(scanner.TOKEN_SEMICOLON, errors.ExpectSemicolonAfterPrintStatement)
	p.EmitByte(byte(common.OpPrint))
}

func (p *Parser) expressionStatement() {
	p.expression()
	p.Consume(scanner.TOKEN_SEMICOLON, errors.ExpectSemicolonAfterExpressionStatement)
	p.EmitByte(byte(common.OpPop))
}

func (p *Parser) synchronize() {
	p.PanicMode = false

	for p.Current.Type != scanner.TOKEN_EOF {
		if p.Previous.Type == scanner.TOKEN_SEMICOLON {
			return
		}

		switch p.Current.Type {
		case scanner.TOKEN_CLASS:
		case scanner.TOKEN_FUN:
		case scanner.TOKEN_VAR:
		case scanner.TOKEN_FOR:
		case scanner.TOKEN_IF:
		case scanner.TOKEN_WHILE:
		case scanner.TOKEN_PRINT:
		case scanner.TOKEN_RETURN:
			return
		}

		p.Advance()
	}
}
