package vm

import (
	"fmt"
	"golox/internal/chunk"
	"golox/pkg/common"
	"golox/pkg/debug"
)

type InterpretResult byte

const (
	InterpretOk InterpretResult = iota
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

		var result InterpretResult

		switch instruction {
		case byte(common.OpConstant), byte(common.OpConstantLong):
			result = v.ConstantOperation(instruction)
			break
		case byte(common.OpAdd):
			result = v.BinaryOperation('+')
			break
		case byte(common.OpSubtract):
			result = v.BinaryOperation('-')
			break
		case byte(common.OpMultiply):
			result = v.BinaryOperation('*')
			break
		case byte(common.OpDivide):
			result = v.BinaryOperation('/')
			break
		case byte(common.OpNegate):
			result = v.UnaryOperation('-')
			break
		case byte(common.OpBang):
			result = v.UnaryOperation('!')
			break
		case byte(common.OpNil):
			v.Stack.Push(common.NilValue())
			break
		case byte(common.OpFalse):
			v.Stack.Push(common.BoolValue(false))
			break
		case byte(common.OpTrue):
			v.Stack.Push(common.BoolValue(true))
			break
		case byte(common.OpReturn):
			if v.DebugMode {
				fmt.Println(debug.PrintValue(v.Stack.Pop()))
			}
			return InterpretOk
		}
		if result != InterpretOk {
			return result
		}
	}
}
