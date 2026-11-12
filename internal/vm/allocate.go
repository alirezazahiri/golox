package vm

import (
	"golox/pkg/bytecode"
	"golox/pkg/common"
)

const (
	initialGCThreshold = 1024 * 1024
	gcHeapGrowFactor   = 2
)

func (v *VM) allocateObject(o common.Obj, size int) {
	v.bytesAllocated += size
	if v.bytesAllocated > v.nextGC {
		v.collectGarbage()
	}

	header := o.Header()
	header.Next = v.objects
	v.objects = o
}

func (v *VM) InternString(chars string) *common.ObjString {
	hash := common.HashString(chars)

	if existing := v.Strings.FindString(chars, hash); existing != nil {
		return existing
	}

	str := &common.ObjString{
		Content: chars,
		Length:  len(chars),
		Hash:    hash,
	}
	v.allocateObject(str, sizeOfString(str))
	v.Strings.Set(str, common.NilValue())
	return str
}

func (v *VM) AllocateFunction() *bytecode.Function {
	fn := bytecode.NewFunction()
	v.allocateObject(fn, sizeOfFunction(fn))
	return fn
}

func (v *VM) allocateClosure(function *bytecode.Function) *bytecode.Closure {
	closure := bytecode.NewClosure(function)
	v.allocateObject(closure, sizeOfClosure(closure))
	return closure
}

func (v *VM) allocateUpvalue(slot int, location *common.Value) *bytecode.Upvalue {
	upvalue := bytecode.NewUpvalue(slot, location)
	v.allocateObject(upvalue, sizeOfUpvalue())
	return upvalue
}

func (v *VM) allocateNative(fn bytecode.NativeFn) *bytecode.ObjNative {
	native := bytecode.NewNative(fn)
	v.allocateObject(native, sizeOfNative())
	return native
}

func sizeOfString(s *common.ObjString) int {
	return 32 + len(s.Content)
}

func sizeOfFunction(_ *bytecode.Function) int {
	return 64
}

func sizeOfClosure(c *bytecode.Closure) int {
	return 40 + c.UpvalueCount*8
}

func sizeOfUpvalue() int {
	return 48
}

func sizeOfNative() int {
	return 24
}

func (v *VM) objectCount() int {
	count := 0
	for obj := v.objects; obj != nil; obj = obj.Header().Next {
		count++
	}
	return count
}
