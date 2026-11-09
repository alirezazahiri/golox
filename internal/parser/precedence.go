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
	prefix := p.getRule(p.Previous.Type)
	if prefix == nil || prefix.Prefix == nil {
		p.error(errors.ExpectExpression)
		return
	}

	canAssign := prec <= PREC_ASSIGNMENT
	prefix.Prefix(canAssign)

	for {
		rule := p.getRule(p.Current.Type)
		if rule == nil || prec > rule.Precedence {
			break
		}
		p.Advance()
		if rule.Infix == nil {
			p.error(errors.ExpectExpression)
			return
		}
		rule.Infix(canAssign)
	}
}
