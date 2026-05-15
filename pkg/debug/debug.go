package debug

import (
	"fmt"
	"golox/internal/chunk"
	"golox/pkg/common"
)

func DisassembleChunk(c *chunk.Chunk, name string) {
	fmt.Printf("== %s ==\n", name)

	for offset := 0; offset < len(c.Code); {
		offset = DisassembleInstruction(c, offset)
	}
}

func DisassembleInstruction(c *chunk.Chunk, offset int) int {
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

	case byte(common.OpConstant):
		return constantInstruction("OP_CONSTANT", c, offset)

	case byte(common.OpConstantLong):
		return constantLongInstruction("OP_CONSTANT_LONG", c, offset)

	case byte(common.OpNegate):
		return simpleInstruction("OP_NEGATE", offset)
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

	case byte(common.OpGreater):
		return simpleInstruction("OP_GREATER", offset)
	case byte(common.OpEqual):
		return simpleInstruction("OP_EQUAL", offset)
	case byte(common.OpLess):
		return simpleInstruction("OP_LESS", offset)

	default:
		fmt.Printf("Unknown opcode %d\n", instruction)
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}

func constantInstruction(name string, chunk *chunk.Chunk, offset int) int {
	constant := chunk.Code[offset+1]
	fmt.Printf("%-16s %4d '%s'\n", name, constant, PrintValue(chunk.Constants.Values[constant]))
	return offset + 2
}

func constantLongInstruction(name string, c *chunk.Chunk, offset int) int {
	b1 := int(c.Code[offset+1])
	b2 := int(c.Code[offset+2])
	b3 := int(c.Code[offset+3])

	constant := (b1 << 16) | (b2 << 8) | b3

	fmt.Printf("%-20s %4d '", name, constant)
	PrintValue(c.Constants.Values[constant])
	fmt.Println("'")

	return offset + 4
}

func PrintValue(value common.Value) string {
	if value.IsNil() {
		return "nil"
	} else if value.IsBool() {
		return fmt.Sprintf("%v", value.As.Bool)
	} else if value.IsString() {
		return fmt.Sprintf("%s", value.AsString().Content)
	}
	return fmt.Sprintf("%v", value.As.Number)
}
