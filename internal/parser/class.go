package parser

import (
	"golox/internal/scanner"
	"golox/pkg/common"
	"golox/pkg/errors"
)

func (p *Parser) classDeclaration() {
	p.Consume(scanner.TOKEN_IDENTIFIER, errors.ExpectClassName)
	className := p.Previous
	nameConstant := p.identifierConstant(&className)
	p.declareVariable()

	p.EmitBytes(byte(common.OpClass), nameConstant)
	p.defineVariable(nameConstant)

	classCompiler := &ClassCompiler{
		Enclosing: p.currentClass,
	}
	p.currentClass = classCompiler

	if p.Match(scanner.TOKEN_LESS) {
		p.Consume(scanner.TOKEN_IDENTIFIER, errors.ExpectSuperclassName)
		p.variable(false)

		if identifiersEqual(&className, &p.Previous) {
			p.error(errors.AClassCantInheritFromItself)
		}

		p.beginScope()
		p.addLocal(syntheticToken("super"))
		p.defineVariable(0)

		p.namedVariable(className, false)
		p.EmitByte(byte(common.OpInherit))
		classCompiler.HasSuperclass = true
	}

	p.namedVariable(className, false)
	p.Consume(scanner.TOKEN_LEFT_BRACE, errors.ExpectLeftBraceBeforeClassBody)
	for !p.Check(scanner.TOKEN_RIGHT_BRACE) && !p.Check(scanner.TOKEN_EOF) {
		p.method()
	}
	p.Consume(scanner.TOKEN_RIGHT_BRACE, errors.ExpectRightBraceAfterBlock)
	p.EmitByte(byte(common.OpPop))

	if classCompiler.HasSuperclass {
		p.endScope()
	}

	p.currentClass = p.currentClass.Enclosing
}

func (p *Parser) method() {
	p.Consume(scanner.TOKEN_IDENTIFIER, errors.ExpectMethodName)
	constant := p.identifierConstant(&p.Previous)

	ty := TypeMethod
	if p.Previous.Lexeme == "init" {
		ty = TypeInitializer
	}
	p.function(ty)
	p.EmitBytes(byte(common.OpMethod), constant)
}

func (p *Parser) dot(canAssign bool) {
	p.Consume(scanner.TOKEN_IDENTIFIER, errors.ExpectPropertyNameAfterDot)
	name := p.identifierConstant(&p.Previous)

	if canAssign && p.Match(scanner.TOKEN_EQUAL) {
		p.expression()
		p.EmitBytes(byte(common.OpSetProperty), name)
	} else if p.Match(scanner.TOKEN_LEFT_PAREN) {
		argCount := p.argumentList()
		p.EmitBytes(byte(common.OpInvoke), name, argCount)
	} else {
		p.EmitBytes(byte(common.OpGetProperty), name)
	}
}

func (p *Parser) this_(_ bool) {
	if p.currentClass == nil {
		p.error(errors.CantUseThisOutsideOfAClass)
		return
	}
	p.variable(false)
}

func (p *Parser) super_(_ bool) {
	if p.currentClass == nil {
		p.error(errors.CantUseSuperOutsideOfAClass)
	} else if !p.currentClass.HasSuperclass {
		p.error(errors.CantUseSuperInAClassWithNoSuperclass)
	}

	p.Consume(scanner.TOKEN_DOT, errors.ExpectSuperclassMethodName)
	p.Consume(scanner.TOKEN_IDENTIFIER, errors.ExpectSuperclassMethodName)
	name := p.identifierConstant(&p.Previous)

	p.namedVariable(syntheticToken("this"), false)
	if p.Match(scanner.TOKEN_LEFT_PAREN) {
		argCount := p.argumentList()
		p.namedVariable(syntheticToken("super"), false)
		p.EmitBytes(byte(common.OpSuperInvoke), name, argCount)
	} else {
		p.namedVariable(syntheticToken("super"), false)
		p.EmitBytes(byte(common.OpGetSuper), name)
	}
}
