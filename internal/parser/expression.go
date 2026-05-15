package parser

import (
	"golox/internal/scanner"
	"golox/pkg/common"
	"golox/pkg/errors"
)

func (p *Parser) Expression() {
	p.parsePrecedence(PREC_ASSIGNMENT)
}

func (p *Parser) grouping() {
	p.Expression()
	p.Consume(scanner.TOKEN_RIGHT_PAREN, errors.ExpectRParenAfterExpression)
}

func (p *Parser) unary() {
	operatorType := p.Previous.Type

	p.parsePrecedence(PREC_UNARY)

	switch operatorType {
	case scanner.TOKEN_MINUS:
		p.EmitByte(byte(common.OpNegate))
		break
	case scanner.TOKEN_BANG:
		p.EmitByte(byte(common.OpBang))
		break
	default:
		return
	}
}

func (p *Parser) binary() {
	// Remember the operator.
	operatorType := p.Previous.Type
	// Compile the right operand.
	rule := p.getRule(operatorType)
	p.parsePrecedence((Precedence)(rule.Precedence + 1))
	// Emit the operator instruction.
	switch operatorType {
	case scanner.TOKEN_PLUS:
		p.EmitByte(byte(common.OpAdd))
		break
	case scanner.TOKEN_MINUS:
		p.EmitByte(byte(common.OpSubtract))
		break
	case scanner.TOKEN_STAR:
		p.EmitByte(byte(common.OpMultiply))
		break
	case scanner.TOKEN_SLASH:
		p.EmitByte(byte(common.OpDivide))
		break
	default:
		return // Unreachable.
	}
}

func (p *Parser) getRule(t scanner.TokenType) *ParserRule {
	rule, exists := p.rules[t]
	if !exists {
		return nil
	}
	return &rule
}
