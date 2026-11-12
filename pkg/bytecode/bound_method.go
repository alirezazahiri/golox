package bytecode

import "golox/pkg/common"

type BoundMethod struct {
	common.Object
	Receiver common.Value
	Method   *Closure
}

func (b *BoundMethod) Type() common.ObjType {
	return common.ObjBoundMethodType
}

func (b *BoundMethod) Header() *common.Object {
	return &b.Object
}

func NewBoundMethod(receiver common.Value, method *Closure) *BoundMethod {
	return &BoundMethod{
		Receiver: receiver,
		Method:   method,
	}
}

func BoundMethodValue(b *BoundMethod) common.Value {
	return common.ObjValue(b)
}

func AsBoundMethod(v common.Value) *BoundMethod {
	return v.As.Obj.(*BoundMethod)
}

func IsBoundMethod(v common.Value) bool {
	return v.IsObj() && v.As.Obj != nil && v.As.Obj.Type() == common.ObjBoundMethodType
}
