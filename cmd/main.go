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
	// v.ActivateDebugMode()

	args := os.Args

	var result vm.InterpretResult

	if len(args) > 1 {
		if len(args) == 2 {
			result = runFile(v, args[1])
		} else {
			fmt.Fprintln(os.Stderr, "Usage: golox [script]")
			os.Exit(64)
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
		fmt.Fprintf(os.Stderr, "Could not read file %q: %v\n", path, fileReadError)
		os.Exit(74)
	}

	return interpreter.Interpret(v, string(bytes))
}

func repl(v *vm.VM) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			return
		}
		line := strings.TrimSpace(scanner.Text())

		if line == "" || line == "exit" {
			return
		}

		result := interpreter.Interpret(v, line)
		if result == vm.InterpretCompileError || result == vm.InterpretRuntimeError {
			// Keep the REPL interactive; errors are already printed.
			continue
		}
	}
}
