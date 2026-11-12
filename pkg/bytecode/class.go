package bytecode

import (
	"golox/pkg/common"
	"golox/pkg/table"
)

type Class struct {
	common.Object
	Name    *common.ObjString
	Methods table.Table
}

func (c *Class) Type() common.ObjType {
	return common.ObjClassType
}

func (c *Class) Header() *common.Object {
	return &c.Object
}

func NewClass(name *common.ObjString) *Class {
	return &Class{
		Name: name,
	}
}

func ClassValue(c *Class) common.Value {
	return common.ObjValue(c)
}

func AsClass(v common.Value) *Class {
	return v.As.Obj.(*Class)
}

func IsClass(v common.Value) bool {
	return v.IsObj() && v.As.Obj != nil && v.As.Obj.Type() == common.ObjClassType
}
