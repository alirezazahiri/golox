package main

import (
	"golox/internal/chunk"
	"golox/internal/vm"
	"golox/pkg/common"
)

func main() {
	v := vm.New()
	v.ActivateDebugMode()
	c := chunk.New()

	c.WriteConstant(1.2, 1) // 1.2
	c.Write(byte(common.OpNegate), 1) // -1.2
	c.Write(byte(common.OpNegate), 1) // 1.2
	c.Write(byte(common.OpNegate), 1) // -1.2
	
	c.WriteConstant(2.8, 1) // 2.8
	c.Write(byte(common.OpNegate), 1) // -2.8
	c.Write(byte(common.OpNegate), 1) // 2.8
	
	c.Write(byte(common.OpSubtract), 1) // -1.2 - -2.8

	c.Write(byte(common.OpReturn), 1)

	v.Interpret(c)
	v.Free()
	c.Free()
}
