package parser

import (
	"golox/internal/scanner"
	"golox/pkg/common"
)

func (p *Parser) literal() {
	switch p.Previous.Type {
	case scanner.TOKEN_FALSE:
		p.EmitByte(byte(common.OpFalse))
		break
	case scanner.TOKEN_NIL:
		p.EmitByte(byte(common.OpNil))
		break
	case scanner.TOKEN_TRUE:
		p.EmitByte(byte(common.OpTrue))
		break
	}
}
