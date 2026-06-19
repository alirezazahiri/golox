package vm

import (
	"fmt"
	"golox/internal/chunk"
	"golox/pkg/common"
	"golox/pkg/debug"
	"golox/pkg/errors"
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
		case byte(common.OpGreater):
			result = v.BinaryOperation('>')
			break
		case byte(common.OpLess):
			result = v.BinaryOperation('<')
			break
		case byte(common.OpEqual):
			result = v.ValuesEqual()
			break
		case byte(common.OpAdd):
			result = v.SumOperation()
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
		case byte(common.OpPrint):
			result = v.PrintValue()
			fmt.Println()
			break
		case byte(common.OpPop):
			v.Stack.Pop()
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
		case byte(common.OpDefineGlobal):
			name, err := v.ReadString()
			if err != nil {
				return InterpretRuntimeError
			}
			value := v.Stack.GetAt(0)
			success := v.Globals.Set(name, value)
			if !success {
				v.runtimeError(errors.FailedToSetValueToVariable, name.Content)
				return InterpretRuntimeError
			}
			v.Stack.Pop()
			break
		case byte(common.OpGetGlobal):
			name, err := v.ReadString()
			if err != nil {
				return InterpretRuntimeError
			}
			value, found := v.Globals.Get(name)
			if !found {
				v.runtimeError(errors.UndefinedVariable, name.Content)
				return InterpretRuntimeError
			}
			v.Stack.Push(value)
			break
		case byte(common.OpSetGlobal):
			name, err := v.ReadString()
			if err != nil {
				return InterpretRuntimeError
			}
			value := v.Stack.GetAt(0)
			v.Globals.Set(name, value)
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
