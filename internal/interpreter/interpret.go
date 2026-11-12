// Package interpreter compiles Lox source to bytecode and runs it on the VM.
package interpreter

import (
	"golox/internal/parser"
	"golox/internal/scanner"
	"golox/internal/vm"
	"golox/pkg/bytecode"
	"golox/pkg/debug"
)

func Interpret(v *vm.VM, source string) vm.InterpretResult {
	function := compile(source, v)
	if function == nil {
		return vm.InterpretCompileError
	}
	return v.InterpretFunction(function)
}

func compile(source string, v *vm.VM) *bytecode.Function {
	s := scanner.New(source)
	p := parser.New(s, v)

	v.SetCompilerRoot(&p.Compiler)
	defer v.SetCompilerRoot(nil)

	p.Advance()

	for !p.Match(scanner.TOKEN_EOF) {
		p.Declaration()
	}

	p.EmitReturn()
	function := p.Compiler.Function

	if !p.HadError && v.DebugMode {
		name := "<script>"
		if function.Name != nil && function.Name.Length != 0 {
			name = function.Name.Content
		}
		debug.DisassembleChunk(function.Chunk, name)
	}

	if p.HadError {
		return nil
	}
	return function
}
