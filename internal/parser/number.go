package parser

import (
	"golox/pkg/common"
	"strconv"
)

func (p *Parser) number(canAssign bool) {
	value, err := strconv.ParseFloat(p.Previous.Lexeme, 64)
	if err != nil {
		return
	}
	p.EmitConstant(common.NumberValue(value))
}
