package bytecode

import "golox/pkg/common"

type Closure struct {
	common.Object
	Function     *Function
	Upvalues     []*Upvalue
	UpvalueCount int
}

func (c *Closure) Type() common.ObjType {
	return common.ObjClosureType
}

func (c *Closure) Header() *common.Object {
	return &c.Object
}

func NewClosure(function *Function) *Closure {
	upvalues := make([]*Upvalue, function.UpvalueCount)
	return &Closure{
		Function:     function,
		Upvalues:     upvalues,
		UpvalueCount: function.UpvalueCount,
	}
}

func ClosureValue(c *Closure) common.Value {
	return common.ObjValue(c)
}

func AsClosure(v common.Value) *Closure {
	return v.As.Obj.(*Closure)
}
