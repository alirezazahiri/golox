package interpreter

import (
	"golox/internal/chunk"
	"golox/internal/parser"
	"golox/internal/scanner"
	"golox/internal/vm"
	"golox/pkg/debug"
)

func Interpret(v *vm.VM, source string) vm.InterpretResult {
	c := chunk.New()

	if !compile(source, c, v) {
		c.Free()
		return vm.InterpretCompileError
	}

	result := v.Interpret(c)

	c.Free()
	return result
}

func compile(source string, c *chunk.Chunk, v *vm.VM) bool {
	s := scanner.New(source)
	p := parser.New(s, c, v)

	p.Advance()
	p.Expression()

	p.EmitReturn()

	if !p.HadError {
		debug.DisassembleChunk(c, "code")
	}
	
	return !p.HadError
}
