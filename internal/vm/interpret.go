package vm

import (
	"fmt"

	"golox/pkg/bytecode"
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

func (v *VM) InterpretFunction(function *bytecode.Function) InterpretResult {
	v.FrameCount = 0
	v.Stack.Reset()

	v.Stack.Push(bytecode.FunctionValue(function))

	frame := &v.Frames[v.FrameCount]
	v.FrameCount++
	frame.Function = function
	frame.IP = 0
	frame.Slots = 0

	return v.Run()
}

func (v *VM) Run() InterpretResult {
	for {
		frame := v.currentFrame()

		if v.DebugMode {
			fmt.Printf("          ")
			for i := range v.Stack.Top() {
				fmt.Printf("[ %s ]", debug.PrintValue(v.Stack.GetAt(i)))
			}
			fmt.Println()
			debug.DisassembleInstruction(frame.Function.Chunk, frame.IP)
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
			value := v.Stack.GetAt(-1)
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
			value := v.Stack.GetAt(-1)
			// Set returns true when the key is new — assignment to an undefined global.
			if v.Globals.Set(name, value) {
				v.Globals.Delete(name)
				v.runtimeError(errors.UndefinedVariable, name.Content)
				return InterpretRuntimeError
			}
			break
		case byte(common.OpGetLocal):
			slot, err := v.ReadByte()
			if err != nil {
				return InterpretRuntimeError
			}
			v.Stack.Push(v.Stack.GetAt(frame.Slots + int(slot)))
			break
		case byte(common.OpSetLocal):
			slot, err := v.ReadByte()
			if err != nil {
				return InterpretRuntimeError
			}
			v.Stack.SetAt(frame.Slots+int(slot), v.Stack.GetAt(-1))
			break
		case byte(common.OpJumpIfFalse):
			offset := v.ReadShort()
			if v.Stack.GetAt(-1).IsFalsey() {
				frame.IP += offset
			}
			break
		case byte(common.OpJump):
			offset := v.ReadShort()
			frame.IP += offset
			break
		case byte(common.OpLoop):
			offset := v.ReadShort()
			frame.IP -= offset
			break
		case byte(common.OpReturn):
			if v.DebugMode {
				fmt.Println(debug.PrintValue(v.Stack.Pop()))
			} else {
				v.Stack.Pop()
			}
			v.FrameCount--
			if v.FrameCount == 0 {
				return InterpretOk
			}
			return InterpretRuntimeError
		}
		if result != InterpretOk {
			return result
		}
	}
}

func (v *VM) ReadShort() int {
	frame := v.currentFrame()
	code := frame.Function.Chunk.Code
	offset := (uint16(code[frame.IP]) << 8) | uint16(code[frame.IP+1])
	frame.IP += 2
	return int(offset)
}
