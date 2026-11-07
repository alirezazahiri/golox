package parser

import "golox/pkg/vm"

func (p *Parser) string(canAssign bool) {
	value := p.Previous.Lexeme[1 : p.Previous.Length-1]
	str := p.vm.InternString(value)
	p.EmitConstant(vm.StringObjValue(str))
}
