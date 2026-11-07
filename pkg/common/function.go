package common

// ObjFunction represents a Lox function object. Chunk holds *bytecode.Chunk at
// runtime; it is typed as any so this package does not depend on pkg/bytecode.
type ObjFunction struct {
	Arity int
	Chunk any
	Name  *ObjString
}

func (f *ObjFunction) Type() ObjType {
	return ObjFunctionType
}

func NewFunction() *ObjFunction {
	return &ObjFunction{
		Arity: 0,
		Name:  nil,
		Chunk: nil,
	}
}
