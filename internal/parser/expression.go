package parser

import (
	"golox/internal/scanner"
	"golox/pkg/common"
	"golox/pkg/errors"
	"math"
)

func (p *Parser) expression() {
	p.parsePrecedence(PREC_ASSIGNMENT)
}

func (p *Parser) grouping(canAssign bool) {
	p.expression()
	p.Consume(scanner.TOKEN_RIGHT_PAREN, errors.ExpectRParenAfterExpression)
}

func (p *Parser) unary(canAssign bool) {
	operatorType := p.Previous.Type

	p.parsePrecedence(PREC_UNARY)

	switch operatorType {
	case scanner.TOKEN_MINUS:
		p.EmitByte(byte(common.OpNegate))
	case scanner.TOKEN_BANG:
		p.EmitByte(byte(common.OpBang))
	default:
		return
	}
}

func (p *Parser) binary(canAssign bool) {
	// Remember the operator.
	operatorType := p.Previous.Type
	// Compile the right operand.
	rule := p.getRule(operatorType)
	p.parsePrecedence((Precedence)(rule.Precedence + 1))
	// Emit the operator instruction.
	switch operatorType {
	case scanner.TOKEN_BANG_EQUAL:
		p.EmitBytes(byte(common.OpEqual), byte(common.OpBang))
	case scanner.TOKEN_EQUAL_EQUAL:
		p.EmitByte(byte(common.OpEqual))
	case scanner.TOKEN_GREATER:
		p.EmitByte(byte(common.OpGreater))
	case scanner.TOKEN_GREATER_EQUAL:
		p.EmitBytes(byte(common.OpLess), byte(common.OpBang))
	case scanner.TOKEN_LESS:
		p.EmitByte(byte(common.OpLess))
	case scanner.TOKEN_LESS_EQUAL:
		p.EmitBytes(byte(common.OpGreater), byte(common.OpBang))
	case scanner.TOKEN_PLUS:
		p.EmitByte(byte(common.OpAdd))
	case scanner.TOKEN_MINUS:
		p.EmitByte(byte(common.OpSubtract))
	case scanner.TOKEN_STAR:
		p.EmitByte(byte(common.OpMultiply))
	case scanner.TOKEN_SLASH:
		p.EmitByte(byte(common.OpDivide))
	default:
		return // Unreachable.
	}
}

func (p *Parser) call(canAssign bool) {
	argCount := p.argumentList()
	p.EmitBytes(byte(common.OpCall), argCount)
}

func (p *Parser) argumentList() byte {
	var argCount byte = 0
	if !p.Check(scanner.TOKEN_RIGHT_PAREN) {
		for {
			p.expression()
			if argCount == math.MaxUint8 {
				p.error(errors.CannotHaveMoreThanMaxCountArgs)
			}
			argCount++

			if !p.Match(scanner.TOKEN_COMMA) {
				break
			}
		}
	}
	p.Consume(scanner.TOKEN_RIGHT_PAREN, errors.ExpectRightParenAfterArgs)
	return argCount
}

func (p *Parser) getRule(t scanner.TokenType) *ParseRule {
	rule, exists := p.rules[t]
	if !exists {
		return nil
	}
	return &rule
}
