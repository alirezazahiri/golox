package vm

import (
	"golox/internal/chunk"
	"golox/pkg/common"
)

type InterpretResult byte

const (
	InterpretOk = iota
	InterpretCompileError
	InterpretRuntimeError
)

func (v *VM) Interpret(chunk *chunk.Chunk) InterpretResult {
	v.Chunk = chunk
	v.IP = 0

	return v.Run()
}

func (v *VM) Run() InterpretResult {
	for {
		instruction, err := v.ReadByte()

		if err != nil {
			return InterpretRuntimeError
		}

		switch instruction {

		case byte(common.OpReturn):
			return InterpretOk
		}
	}
}

func (v *VM) ReadByte() (byte, error) {
	b := v.Chunk.Code[v.IP]
	v.IP++
	return b, nil
}
