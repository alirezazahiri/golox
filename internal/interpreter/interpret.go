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
	function := compile(source, c, v)

	if function == nil {
		c.Free()
		return vm.InterpretCompileError
	}

	result := v.Interpret(c)

	c.Free()
	return result
}

func compile(source string, c *bytecode.Chunk, v *vm.VM) *bytecode.Function {
	s := scanner.New(source)
	p := parser.New(s, c, v)

	p.Advance()

	for !p.Match(scanner.TOKEN_EOF) {
		p.Declaration()
	}

	p.EmitReturn()
	function := p.Compiler.Function

	if !p.HadError {
		var name string
		if function.Name != nil && function.Name.Length != 0 {
			name = function.Name.Content
		} else {
			name = "<script>"
		}
		debug.DisassembleChunk(c, name)
	}

	return function
}
