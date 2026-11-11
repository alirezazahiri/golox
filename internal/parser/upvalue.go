package parser

import (
	"golox/internal/scanner"
	"golox/pkg/errors"
)

func (p *Parser) addUpvalue(compiler *Compiler, index uint8, isLocal bool) int {
	for i := 0; i < compiler.UpvalueCount; i++ {
		upvalue := &compiler.Upvalues[i]
		if upvalue.Index == index && upvalue.IsLocal == isLocal {
			return i
		}
	}

	if compiler.UpvalueCount == uInt8Count {
		p.error(errors.TooManyClosureVariables)
		return 0
	}

	compiler.Upvalues[compiler.UpvalueCount].Index = index
	compiler.Upvalues[compiler.UpvalueCount].IsLocal = isLocal
	compiler.UpvalueCount++
	compiler.Function.UpvalueCount = compiler.UpvalueCount
	return compiler.UpvalueCount - 1
}

func (p *Parser) resolveUpvalue(compiler *Compiler, name scanner.Token) int {
	if compiler.Enclosing == nil {
		return -1
	}

	if local, ok := p.resolveLocalIn(compiler.Enclosing, name); ok {
		compiler.Enclosing.Locals[local].IsCaptured = true
		return p.addUpvalue(compiler, uint8(local), true)
	}

	upvalue := p.resolveUpvalue(compiler.Enclosing, name)
	if upvalue != -1 {
		return p.addUpvalue(compiler, uint8(upvalue), false)
	}

	return -1
}

func (p *Parser) resolveLocalIn(compiler *Compiler, name scanner.Token) (int, bool) {
	for i := compiler.LocalCount - 1; i >= 0; i-- {
		local := &compiler.Locals[i]
		if identifiersEqual(&local.Name, &name) {
			if local.Depth == -1 {
				p.error(errors.CannotAccessVariableInItsInitializer)
			}
			return i, true
		}
	}
	return 0, false
}
