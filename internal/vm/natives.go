package vm

import (
	"golox/pkg/bytecode"
	"golox/pkg/common"
	"time"
)

func (v *VM) DefineNative(name string, fn bytecode.NativeFn) {
	nameObj := v.InternString(name)
	native := bytecode.NewNative(fn)
	v.Globals.Set(nameObj, bytecode.NativeFnValue(native))
}

// NOTE: native fns
// these methods MUST adhere to the following signature:
// type NativeFn func(argCount int, args []common.Value) common.Value
func clockNative(_ int, _ []common.Value) common.Value {
	return common.NumberValue(float64(time.Now().UnixMilli()))
}

// NOTE: define natives here
func (v *VM) defineNatives() {
	v.DefineNative("clock", clockNative)
}
