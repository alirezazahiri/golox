package parser

import "golox/pkg/common"

func (p *Parser) string() {
	value := p.Previous.Lexeme[1 : p.Previous.Length-1]
	str := p.vm.InternString(value)
	p.EmitConstant(common.StringObjValue(str))
}
