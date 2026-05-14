package vm

import (
	"fmt"
	"golox/internal/chunk"
	"golox/pkg/common"
	"golox/pkg/debug"
)

type InterpretResult byte

const (
	InterpretOk = iota
	InterpretCompileError
	InterpretRuntimeError
)

func (v *VM) Interpret(chunk *chunk.Chunk) InterpretResult {
	v.WriteChunk(chunk)
	v.IP = 0
	return v.Run()
}

func (v *VM) Run() InterpretResult {
	for {
		if v.DebugMode {
			fmt.Printf("          ")
			for i := range v.Stack.Top() {
				fmt.Printf("[ %s ]", debug.PrintValue(v.Stack.GetAt(i)))
			}
			fmt.Println()
			debug.DisassembleInstruction(v.Chunk, v.IP)
		}

		instruction, err := v.ReadByte()

		if err != nil {
			return InterpretRuntimeError
		}

		switch instruction {
		case byte(common.OpConstant):
			c, err := v.ReadConstant()
			if err != nil {
				return InterpretRuntimeError
			}
			v.Stack.Push(c)
			fmt.Println(debug.PrintValue(c))
			break
		case byte(common.OpConstantLong):
			c, err := v.ReadConstantLong()
			if err != nil {
				return InterpretRuntimeError
			}
			v.Stack.Push(c)
			fmt.Println(debug.PrintValue(c))
			break
		case byte(common.OpAdd):
			v.BinaryOperation('+')
			break
		case byte(common.OpSubtract):
			v.BinaryOperation('-')
			break
		case byte(common.OpMultiply):
			v.BinaryOperation('*')
			break
		case byte(common.OpDivide):
			v.BinaryOperation('/')
			break
		case byte(common.OpNegate):
			v.Stack.UpdateTop(func(v common.Value) common.Value { return -v })
			break
		case byte(common.OpReturn):
			fmt.Println(debug.PrintValue(v.Stack.Pop()))
			return InterpretOk
		}
	}
}
