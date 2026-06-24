package parser

import (
	"golox/internal/scanner"
	"golox/pkg/common"
	"golox/pkg/errors"
)

func (p *Parser) parseVariable(errorMessage string) uint8 {
	p.Consume(scanner.TOKEN_IDENTIFIER, errorMessage)
	p.declareVariable()
	if p.compiler.ScopeDepth > 0 {
		return 0
	}
	return p.identifierConstant(&p.Previous)
}

func (p *Parser) identifierConstant(token *scanner.Token) uint8 {
	str := p.vm.InternString(token.Lexeme)
	return uint8(p.chunk.AddConstant(common.StringObjValue(str)))
}

func (p *Parser) defineVariable(name uint8) {
	if p.compiler.ScopeDepth > 0 {
		p.compiler.markInitialized()
		return
	}
	p.EmitBytes(byte(common.OpDefineGlobal), name)
}

func (p *Parser) declareVariable() {
	if p.compiler.ScopeDepth == 0 {
		return
	}
	name := p.Previous
	for i := p.compiler.LocalCount - 1; i >= 0; i-- {
		local := p.compiler.Locals[i]
		if local.Depth != -1 && local.Depth < p.compiler.ScopeDepth {
			break
		}
		if identifiersEqual(&name, &local.Name) {
			p.error(errors.AlreadyVariableWithThisNameInThisScope)
		}
	}
	p.compiler.addLocal(name)
}

func identifiersEqual(a *scanner.Token, b *scanner.Token) bool {
	return a.Lexeme == b.Lexeme
}

func (p *Parser) variable(canAssign bool) {
	p.namedVariable(p.Previous, canAssign)
}

func (p *Parser) namedVariable(name scanner.Token, canAssign bool) {
	var arg uint8
	var getOp uint8
	var setOp uint8

	if local, ok := p.compiler.resolveLocal(name); ok {
		arg = uint8(local)
		getOp = uint8(common.OpGetLocal)
		setOp = uint8(common.OpSetLocal)
	} else {
		arg = p.identifierConstant(&name)
		getOp = uint8(common.OpGetGlobal)
		setOp = uint8(common.OpSetGlobal)
	}

	if canAssign && p.Match(scanner.TOKEN_EQUAL) {
		p.expression()
		p.EmitBytes(setOp, arg)
	} else {
		p.EmitBytes(getOp, arg)
	}
}
