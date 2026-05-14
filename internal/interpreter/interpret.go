package interpreter

import (
	"golox/internal/scanner"
	"golox/internal/vm"
)

func Interpret(source string) vm.InterpretResult {
	compile(source)
	return vm.InterpretOk
}

func compile(source string) {
	s := scanner.New(source)
	s.ScanTokens()
}