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
	v.OpenUpvalues = nil

	closure := bytecode.NewClosure(function)
	v.Stack.Push(bytecode.ClosureValue(closure))

	v.callClosure(closure, 0)
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
			debug.DisassembleInstruction(frame.Closure.Function.Chunk, frame.IP)
		}

		instruction, err := v.ReadByte()

		if err != nil {
			return InterpretRuntimeError
		}

		var result InterpretResult

		switch instruction {
		case byte(common.OpConstant), byte(common.OpConstantLong):
			result = v.ConstantOperation(instruction)
		case byte(common.OpGreater):
			result = v.BinaryOperation('>')
		case byte(common.OpLess):
			result = v.BinaryOperation('<')
		case byte(common.OpEqual):
			result = v.ValuesEqual()
		case byte(common.OpAdd):
			result = v.SumOperation()
		case byte(common.OpSubtract):
			result = v.BinaryOperation('-')
		case byte(common.OpMultiply):
			result = v.BinaryOperation('*')
		case byte(common.OpDivide):
			result = v.BinaryOperation('/')
		case byte(common.OpNegate):
			result = v.UnaryOperation('-')
		case byte(common.OpPrint):
			result = v.PrintValue()
			fmt.Println()
		case byte(common.OpPop):
			v.Stack.Pop()
		case byte(common.OpCloseUpvalue):
			v.closeUpvalues(v.Stack.Top() - 1)
			v.Stack.Pop()
		case byte(common.OpBang):
			result = v.UnaryOperation('!')
		case byte(common.OpNil):
			v.Stack.Push(common.NilValue())
		case byte(common.OpFalse):
			v.Stack.Push(common.BoolValue(false))
		case byte(common.OpTrue):
			v.Stack.Push(common.BoolValue(true))
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
		case byte(common.OpGetLocal):
			slot, err := v.ReadByte()
			if err != nil {
				return InterpretRuntimeError
			}
			v.Stack.Push(v.Stack.GetAt(frame.Slots + int(slot)))
		case byte(common.OpSetLocal):
			slot, err := v.ReadByte()
			if err != nil {
				return InterpretRuntimeError
			}
			v.Stack.SetAt(frame.Slots+int(slot), v.Stack.GetAt(-1))
		case byte(common.OpGetUpvalue):
			slot, err := v.ReadByte()
			if err != nil {
				return InterpretRuntimeError
			}
			v.Stack.Push(*frame.Closure.Upvalues[slot].Location)
		case byte(common.OpSetUpvalue):
			slot, err := v.ReadByte()
			if err != nil {
				return InterpretRuntimeError
			}
			*frame.Closure.Upvalues[slot].Location = v.Stack.GetAt(-1)
		case byte(common.OpJumpIfFalse):
			offset := v.ReadShort()
			if v.Stack.GetAt(-1).IsFalsey() {
				frame.IP += offset
			}
		case byte(common.OpJump):
			offset := v.ReadShort()
			frame.IP += offset
		case byte(common.OpLoop):
			offset := v.ReadShort()
			frame.IP -= offset
		case byte(common.OpClosure):
			functionValue, err := v.ReadConstant()
			if err != nil {
				return InterpretRuntimeError
			}
			function := bytecode.AsFunction(functionValue)
			closure := bytecode.NewClosure(function)
			for i := 0; i < closure.UpvalueCount; i++ {
				isLocal, err := v.ReadByte()
				if err != nil {
					return InterpretRuntimeError
				}
				index, err := v.ReadByte()
				if err != nil {
					return InterpretRuntimeError
				}
				if isLocal == 1 {
					closure.Upvalues[i] = v.captureUpvalue(frame.Slots + int(index))
				} else {
					closure.Upvalues[i] = frame.Closure.Upvalues[index]
				}
			}
			v.Stack.Push(bytecode.ClosureValue(closure))
		case byte(common.OpReturn):
			result := v.Stack.Pop()
			if v.DebugMode {
				fmt.Println(debug.PrintValue(result))
			}
			slots := frame.Slots
			v.closeUpvalues(slots)
			v.FrameCount--
			if v.FrameCount == 0 {
				if v.Stack.Top() > 0 {
					v.Stack.Pop()
				}
				return InterpretOk
			}
			v.Stack.TruncateTo(slots)
			v.Stack.Push(result)
			frame = &v.Frames[v.FrameCount-1]
		case byte(common.OpCall):
			argCount, err := v.ReadByte()
			if err != nil {
				return InterpretRuntimeError
			}
			if !v.callValue(v.Stack.GetAt(-(int(argCount) + 1)), int(argCount)) {
				return InterpretRuntimeError
			}
			frame = &v.Frames[v.FrameCount-1]
		}
		if result != InterpretOk {
			return result
		}
	}
}

func (v *VM) ReadShort() int {
	frame := v.currentFrame()
	code := frame.Closure.Function.Chunk.Code
	offset := (uint16(code[frame.IP]) << 8) | uint16(code[frame.IP+1])
	frame.IP += 2
	return int(offset)
}

func (v *VM) callValue(callee common.Value, argCount int) bool {
	if callee.IsObj() {
		switch callee.As.Obj.Type() {
		case common.ObjClosureType:
			return v.callClosure(bytecode.AsClosure(callee), argCount)
		case common.ObjNativeFnType:
			native := bytecode.AsNative(callee)
			args := v.Stack.TopSlice(argCount)
			result := native.Function(argCount, args)
			v.Stack.TruncateTo(v.Stack.Top() - argCount - 1)
			v.Stack.Push(result)
			return true
		}
	}
	v.runtimeError(errors.CanOnlyCallFunctionsAndClasses)
	return false
}

func (v *VM) callClosure(closure *bytecode.Closure, argCount int) bool {
	if argCount != closure.Function.Arity {
		v.runtimeError(errors.ExpectedArgumentsButGot, closure.Function.Arity, argCount)
		return false
	}
	if v.FrameCount == FramesMax {
		v.runtimeError(errors.StackOverflow)
		return false
	}
	frame := &v.Frames[v.FrameCount]
	v.FrameCount++
	frame.Closure = closure
	frame.IP = 0
	frame.Slots = v.Stack.Top() - argCount - 1
	return true
}
