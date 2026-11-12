package vm

import (
	"golox/pkg/bytecode"
	"golox/pkg/common"
	"golox/pkg/errors"
)

func (v *VM) bindMethod(klass *bytecode.Class, name *common.ObjString) bool {
	method, found := klass.Methods.Get(name)
	if !found {
		v.runtimeError(errors.UndefinedProperty, name.Content)
		return false
	}
	bound := v.allocateBoundMethod(v.Stack.GetAt(-1), bytecode.AsClosure(method))
	v.Stack.Pop()
	v.Stack.Push(bytecode.BoundMethodValue(bound))
	return true
}

func (v *VM) invokeFromClass(klass *bytecode.Class, name *common.ObjString, argCount int) bool {
	method, found := klass.Methods.Get(name)
	if !found {
		v.runtimeError(errors.UndefinedMethod, name.Content)
		return false
	}
	return v.callClosure(bytecode.AsClosure(method), argCount)
}

func (v *VM) invoke(name *common.ObjString, argCount int) bool {
	receiver := v.Stack.GetAt(-(argCount + 1))
	if !bytecode.IsInstance(receiver) {
		v.runtimeError(errors.OnlyInstancesHaveMethods)
		return false
	}
	instance := bytecode.AsInstance(receiver)

	if value, found := instance.Fields.Get(name); found {
		v.Stack.SetAt(v.Stack.Top()-(argCount+1), value)
		return v.callValue(value, argCount)
	}
	return v.invokeFromClass(instance.Klass, name, argCount)
}

func (v *VM) defineMethod(name *common.ObjString) {
	method := v.Stack.GetAt(-1)
	klass := bytecode.AsClass(v.Stack.GetAt(-2))
	klass.Methods.Set(name, method)
	v.Stack.Pop()
}
