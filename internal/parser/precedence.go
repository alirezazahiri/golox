package parser

import "golox/pkg/errors"

type Precedence uint8

const (
	PREC_NONE       Precedence = iota
	PREC_ASSIGNMENT            // =
	PREC_OR                    // or
	PREC_AND                   // and
	PREC_EQUALITY              // == !=
	PREC_COMPARISON            // < > <= >=
	PREC_TERM                  // + -
	PREC_FACTOR                // * /
	PREC_UNARY                 // ! -
	PREC_CALL                  // . ()
	PREC_PRIMARY
)

func (p *Parser) parsePrecedence(prec Precedence) {
	p.Advance()
	prefixRule := p.getRule(p.Previous.Type).Prefix

	if prefixRule == nil {
		p.error(errors.ExpectExpression)
		return
	}
	
	prefixRule()
	
	for prec <= p.getRule(p.Current.Type).Precedence {
		p.Advance()
		infixRule := p.getRule(p.Previous.Type).Infix
		infixRule()
	}
}
