package parser

import (
	"golox/internal/scanner"
	"golox/pkg/common"
)

func (p *Parser) parseVariable(errorMessage string) uint8 {
	p.Consume(scanner.TOKEN_IDENTIFIER, errorMessage)
	return p.identifierConstant(&p.Previous)
}

func (p *Parser) identifierConstant(token *scanner.Token) uint8 {
	str := p.vm.InternString(token.Lexeme)
	return uint8(p.chunk.AddConstant(common.StringObjValue(str)))
}

func (p *Parser) defineVariable(name uint8) {
	p.EmitBytes(byte(common.OpDefineGlobal), name)
}

func (p *Parser) variable() {
	p.namedVariable(p.Previous)
}

func (p *Parser) namedVariable(name scanner.Token) {
	arg := p.identifierConstant(&name)
	p.EmitBytes(byte(common.OpGetGlobal), arg)
}