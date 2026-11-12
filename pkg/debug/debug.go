package debug

import (
	"fmt"

	"golox/pkg/bytecode"
	"golox/pkg/common"
)

func DisassembleChunk(c *bytecode.Chunk, name string) {
	fmt.Printf("== %s ==\n", name)

	for offset := 0; offset < len(c.Code); {
		offset = DisassembleInstruction(c, offset)
	}
}

func DisassembleInstruction(c *bytecode.Chunk, offset int) int {
	fmt.Printf("%04d ", offset)

	line := c.GetLine(offset)

	if offset > 0 && line == c.GetLine(offset-1) {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", line)
	}

	instruction := c.Code[offset]

	switch instruction {

	case byte(common.OpReturn):
		return simpleInstruction("OP_RETURN", offset)
	case byte(common.OpCall):
		return constantInstruction("OP_CALL", c, offset)

	case byte(common.OpConstant):
		return constantInstruction("OP_CONSTANT", c, offset)

	case byte(common.OpConstantLong):
		return constantLongInstruction("OP_CONSTANT_LONG", c, offset)

	case byte(common.OpNegate):
		return simpleInstruction("OP_NEGATE", offset)
	case byte(common.OpPrint):
		return simpleInstruction("OP_PRINT", offset)
	case byte(common.OpPop):
		return simpleInstruction("OP_POP", offset)
	case byte(common.OpCloseUpvalue):
		return simpleInstruction("OP_CLOSE_UPVALUE", offset)
	case byte(common.OpBang):
		return simpleInstruction("OP_BANG", offset)

	case byte(common.OpAdd):
		return simpleInstruction("OP_ADD", offset)
	case byte(common.OpSubtract):
		return simpleInstruction("OP_SUBTRACT", offset)
	case byte(common.OpMultiply):
		return simpleInstruction("OP_MULTIPLY", offset)
	case byte(common.OpDivide):
		return simpleInstruction("OP_DIVIDE", offset)

	case byte(common.OpNil):
		return simpleInstruction("OP_NIL", offset)
	case byte(common.OpFalse):
		return simpleInstruction("OP_FALSE", offset)
	case byte(common.OpTrue):
		return simpleInstruction("OP_TRUE", offset)

	case byte(common.OpDefineGlobal):
		return constantInstruction("OP_DEFINE_GLOBAL", c, offset)
	case byte(common.OpGetGlobal):
		return constantInstruction("OP_GET_GLOBAL", c, offset)
	case byte(common.OpSetGlobal):
		return constantInstruction("OP_SET_GLOBAL", c, offset)
	case byte(common.OpGetLocal):
		return byteInstruction("OP_GET_LOCAL", c, offset)
	case byte(common.OpSetLocal):
		return byteInstruction("OP_SET_LOCAL", c, offset)
	case byte(common.OpGetUpvalue):
		return byteInstruction("OP_GET_UPVALUE", c, offset)
	case byte(common.OpSetUpvalue):
		return byteInstruction("OP_SET_UPVALUE", c, offset)
	case byte(common.OpClosure):
		return closureInstruction("OP_CLOSURE", c, offset)
	case byte(common.OpGreater):
		return simpleInstruction("OP_GREATER", offset)
	case byte(common.OpEqual):
		return simpleInstruction("OP_EQUAL", offset)
	case byte(common.OpLess):
		return simpleInstruction("OP_LESS", offset)

	case byte(common.OpJump):
		return jumpInstruction("OP_JUMP", 1, c, offset)
	case byte(common.OpJumpIfFalse):
		return jumpInstruction("OP_JUMP_IF_FALSE", 1, c, offset)
	case byte(common.OpLoop):
		return jumpInstruction("OP_LOOP", -1, c, offset)

	case byte(common.OpClass):
		return constantInstruction("OP_CLASS", c, offset)
	case byte(common.OpGetProperty):
		return constantInstruction("OP_GET_PROPERTY", c, offset)
	case byte(common.OpSetProperty):
		return constantInstruction("OP_SET_PROPERTY", c, offset)
	case byte(common.OpMethod):
		return constantInstruction("OP_METHOD", c, offset)
	case byte(common.OpInvoke):
		return invokeInstruction("OP_INVOKE", c, offset)
	case byte(common.OpInherit):
		return simpleInstruction("OP_INHERIT", offset)
	case byte(common.OpGetSuper):
		return constantInstruction("OP_GET_SUPER", c, offset)
	case byte(common.OpSuperInvoke):
		return invokeInstruction("OP_SUPER_INVOKE", c, offset)

	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}

func constantInstruction(name string, chunk *bytecode.Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '%s'\n", name, constant, PrintValue(chunk.Constants.Values[constant]))
	return offset + 2
}

func constantLongInstruction(name string, c *bytecode.Chunk, offset int) int {
	b1 := int(c.Code[offset+1])
	b2 := int(c.Code[offset+2])
	b3 := int(c.Code[offset+3])

	constant := (b1 << 16) | (b2 << 8) | b3

	fmt.Printf("%-20s %4d '", name, constant)
	PrintValue(c.Constants.Values[constant])
	fmt.Println("'")

	return offset + 4
}

func byteInstruction(name string, c *bytecode.Chunk, offset int) int {
	slot := c.Code[offset+1]
	fmt.Printf("%-16s %4d\n", name, slot)
	return offset + 2
}

func closureInstruction(name string, c *bytecode.Chunk, offset int) int {
	constant := c.Code[offset+1]
	fmt.Printf("%-16s %4d ", name, constant)
	fmt.Printf("%s\n", PrintValue(c.Constants.Values[constant]))
	offset += 2

	function := bytecode.AsFunction(c.Constants.Values[constant])
	for j := 0; j < function.UpvalueCount; j++ {
		isLocal := c.Code[offset]
		index := c.Code[offset+1]
		fmt.Printf("%04d      |                     %s %d\n",
			offset, map[bool]string{true: "local", false: "upvalue"}[isLocal == 1], index)
		offset += 2
	}
	return offset
}

func jumpInstruction(name string, sign int, c *bytecode.Chunk, offset int) int {
	jump := uint16(c.Code[offset+1]) << 8
	jump |= uint16(c.Code[offset+2])

	fmt.Printf("%-16s %4d -> %d\n", name, offset, offset+3+sign*int(jump))
	return offset + 3
}

func invokeInstruction(name string, c *bytecode.Chunk, offset int) int {
	constant := c.Code[offset+1]
	argCount := c.Code[offset+2]
	fmt.Printf("%-16s (%d args) %4d '%s'\n",
		name, argCount, constant, PrintValue(c.Constants.Values[constant]))
	return offset + 3
}

func PrintValue(value common.Value) string {
	if value.IsNil() {
		return "nil"
	} else if value.IsBool() {
		return fmt.Sprintf("%v", value.As.Bool)
	} else if value.IsString() {
		return value.AsString().Content
	} else if value.IsFunction() {
		if fn := bytecode.AsFunction(value); fn.Name != nil {
			return fmt.Sprintf("<fn %s>", fn.Name.Content)
		}
		return "<script>"
	} else if value.IsClosure() {
		if fn := bytecode.AsClosure(value).Function; fn.Name != nil {
			return fmt.Sprintf("<fn %s>", fn.Name.Content)
		}
		return "<script>"
	} else if value.IsNativeFn() {
		return "<native fn>"
	} else if value.IsClass() {
		return fmt.Sprintf("<class %s>", bytecode.AsClass(value).Name.Content)
	} else if value.IsInstance() {
		return fmt.Sprintf("%s instance", bytecode.AsInstance(value).Klass.Name.Content)
	} else if value.IsBoundMethod() {
		fn := bytecode.AsBoundMethod(value).Method.Function
		if fn.Name != nil {
			return fmt.Sprintf("<fn %s>", fn.Name.Content)
		}
		return "<fn>"
	}
	return fmt.Sprintf("%v", value.As.Number)
}
