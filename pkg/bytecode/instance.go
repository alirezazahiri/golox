package bytecode

import (
	"golox/pkg/common"
	"golox/pkg/table"
)

type Instance struct {
	common.Object
	Klass  *Class
	Fields table.Table
}

func (i *Instance) Type() common.ObjType {
	return common.ObjInstanceType
}

func (i *Instance) Header() *common.Object {
	return &i.Object
}

func NewInstance(klass *Class) *Instance {
	return &Instance{
		Klass: klass,
	}
}

func InstanceValue(i *Instance) common.Value {
	return common.ObjValue(i)
}

func AsInstance(v common.Value) *Instance {
	return v.As.Obj.(*Instance)
}

func IsInstance(v common.Value) bool {
	return v.IsObj() && v.As.Obj != nil && v.As.Obj.Type() == common.ObjInstanceType
}
