package bytecode

import "golox/pkg/common"

func NewFunctionObject() *common.ObjFunction {
	fn := common.NewFunction()
	fn.Chunk = NewChunk()
	return fn
}

func FunctionChunk(fn *common.ObjFunction) *Chunk {
	if fn == nil || fn.Chunk == nil {
		return nil
	}
	return fn.Chunk.(*Chunk)
}

func SetFunctionChunk(fn *common.ObjFunction, c *Chunk) {
	fn.Chunk = c
}
