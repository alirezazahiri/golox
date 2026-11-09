package parser

import (
	"strconv"

	"golox/pkg/common"
	"golox/pkg/errors"
)

func (p *Parser) number(canAssign bool) {
	value, err := strconv.ParseFloat(p.Previous.Lexeme, 64)
	if err != nil {
		p.error(errors.InvalidNumber)
		return
	}
	p.EmitConstant(common.NumberValue(value))
}
