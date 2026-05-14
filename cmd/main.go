package main

import (
	"golox/internal/chunk"
	"golox/internal/vm"
	"golox/pkg/common"
	"golox/pkg/debug"
)

func main() {
	v := vm.New()
	c := chunk.New()

	constant := c.AddConstant(1.2)
	c.Write(byte(common.OpConstant), 1)
	c.Write(byte(common.OpCode(constant)), 1)

	c.Write(byte(common.OpReturn), 1)
	c.Write(byte(common.OpReturn), 1)
	c.Write(byte(common.OpReturn), 1)
	c.Write(byte(common.OpReturn), 1)
	c.Write(byte(common.OpReturn), 2)
	c.Write(byte(common.OpReturn), 2)
	c.Write(byte(common.OpReturn), 5)

	debug.DisassembleChunk(c, "test chunk")
	v.Free()
	c.Free()
}
