package main

import (
	"bufio"
	"fmt"
	"golox/internal/interpreter"
	"golox/internal/vm"
	"golox/pkg/file"
	"os"
	"strings"
)

func main() {
	v := vm.New()
	v.ActivateDebugMode()

	args := os.Args

	files := make([]string, len(args))

	var result vm.InterpretResult

	if len(args) > 1 {
		files = args[1:]
		if files != nil && len(files) == 1 {
			result = runFile(v, files[0])
		} else if len(files) > 1 {
			panic("Usage: golox [script]")
		}
	} else {
		repl(v)
	}

	v.Free()

	if result == vm.InterpretCompileError {
		os.Exit(65)
	}
	if result == vm.InterpretRuntimeError {
		os.Exit(70)
	}
}

func runFile(v *vm.VM, path string) vm.InterpretResult {
	f := file.New(path)

	bytes, fileReadError := f.Read()

	if fileReadError != nil {
		panic(fileReadError)
	}

	fileContent := string(bytes)

	return interpreter.Interpret(v, fileContent)
}

func repl(v *vm.VM) {
	for {
		fmt.Print("> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())

		if line == "" || line == "exit" {
			return
		}

		interpreter.Interpret(v, line)
	}
}

// c := chunk.New()

// c.WriteConstant(1.2, 1) // 1.2
// c.Write(byte(common.OpNegate), 1) // -1.2
// c.Write(byte(common.OpNegate), 1) // 1.2
// c.Write(byte(common.OpNegate), 1) // -1.2

// c.WriteConstant(2.8, 1) // 2.8
// c.Write(byte(common.OpNegate), 1) // -2.8
// c.Write(byte(common.OpNegate), 1) // 2.8

// c.Write(byte(common.OpSubtract), 1) // -1.2 - -2.8

// c.Write(byte(common.OpReturn), 1)
// v.Interpret(c)
