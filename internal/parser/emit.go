package parser

import (
	"golox/pkg/common"
	"golox/pkg/errors"
)

func (p *Parser) EmitBytes(bs ...byte) {
	for _, b := range bs {
		p.chunk.Write(b, p.scanner.Line)
	}
}

func (p *Parser) EmitByte(b byte) {
	p.EmitBytes(b)
}

func (p *Parser) EmitReturn() {
	p.EmitByte(byte(common.OpReturn))
}

func (p *Parser) EmitConstant(value common.Value) {
	constant := p.chunk.AddConstant(value)

	if constant <= 0xff {
		p.EmitBytes(
			byte(common.OpConstant),
			byte(constant),
		)
	} else if constant <= 0xFFFFFF {
		p.EmitByte(byte(common.OpConstantLong))

		p.EmitByte(byte((constant >> 16) & 0xFF))
		p.EmitByte(byte((constant >> 8) & 0xFF))
		p.EmitByte(byte(constant & 0xFF))
	} else {
		p.error(errors.TooManyConstantsInOneChunk)
	}
}
