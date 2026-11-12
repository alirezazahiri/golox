package bytecode

import "golox/pkg/common"

type Function struct {
	common.Object
	Arity         int
	UpvalueCount  int
	Chunk         *Chunk
	Name          *common.ObjString
	IsInitializer bool
}

func (f *Function) Type() common.ObjType {
	return common.ObjFunctionType
}

func (f *Function) Header() *common.Object {
	return &f.Object
}

func NewFunction() *Function {
	return &Function{
		Arity: 0,
		Name:  nil,
		Chunk: NewChunk(),
	}
}

func FunctionValue(fn *Function) common.Value {
	return common.ObjValue(fn)
}

func AsFunction(v common.Value) *Function {
	return v.As.Obj.(*Function)
}
