package interpreter

import (
	"golox/internal/parser"
	"golox/internal/scanner"
	"golox/internal/vm"
	"golox/pkg/bytecode"
	"golox/pkg/debug"
)

func Interpret(v *vm.VM, source string) vm.InterpretResult {
	c := bytecode.NewChunk()

	if !compile(source, c, v) {
		c.Free()
		return vm.InterpretCompileError
	}

	result := v.Interpret(c)

	c.Free()
	return result
}

func compile(source string, c *bytecode.Chunk, v *vm.VM) bool {
	s := scanner.New(source)
	p := parser.New(s, c, v)

	p.Advance()

	for !p.Match(scanner.TOKEN_EOF) {
		p.Declaration()
	}

	p.EmitReturn()

	if !p.HadError {
		debug.DisassembleChunk(c, "code")
	}

	return !p.HadError
}
