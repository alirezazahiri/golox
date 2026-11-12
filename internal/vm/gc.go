package vm

import (
	"golox/internal/parser"
	"golox/pkg/bytecode"
	"golox/pkg/common"
	"golox/pkg/table"
)

func (v *VM) SetCompilerRoot(compiler **parser.Compiler) {
	v.compilerRoot = compiler
}

func (v *VM) CollectGarbage() {
	v.collectGarbage()
}

func (v *VM) collectGarbage() {
	v.markRoots()
	v.traceReferences()
	v.Strings.RemoveWhite()
	v.sweep()

	v.nextGC = v.bytesAllocated * gcHeapGrowFactor
	if v.nextGC < initialGCThreshold {
		v.nextGC = initialGCThreshold
	}
}

func (v *VM) markRoots() {
	for i := 0; i < v.Stack.Top(); i++ {
		v.markValue(v.Stack.GetAt(i))
	}

	for i := 0; i < v.FrameCount; i++ {
		v.markObject(v.Frames[i].Closure)
	}

	for upvalue := v.OpenUpvalues; upvalue != nil; upvalue = upvalue.Next {
		v.markObject(upvalue)
	}

	v.markTable(v.Globals)
	v.markObject(v.InitString)
	v.markCompilerRoots()
}

func (v *VM) markCompilerRoots() {
	if v.compilerRoot == nil || *v.compilerRoot == nil {
		return
	}
	for compiler := *v.compilerRoot; compiler != nil; compiler = compiler.Enclosing {
		v.markObject(compiler.Function)
	}
}

func (v *VM) markTable(t *table.Table) {
	for i := range t.Entries {
		e := &t.Entries[i]
		if e.Key == nil {
			continue
		}
		v.markObject(e.Key)
		v.markValue(e.Value)
	}
}

func (v *VM) markValue(value common.Value) {
	if value.IsObj() {
		v.markObject(value.As.Obj)
	}
}

func (v *VM) markObject(obj common.Obj) {
	if obj == nil {
		return
	}
	header := obj.Header()
	if header.Marked {
		return
	}
	header.Marked = true
	v.grayStack = append(v.grayStack, obj)
}

func (v *VM) traceReferences() {
	for len(v.grayStack) > 0 {
		obj := v.grayStack[len(v.grayStack)-1]
		v.grayStack = v.grayStack[:len(v.grayStack)-1]
		v.blackenObject(obj)
	}
}

func (v *VM) blackenObject(obj common.Obj) {
	switch o := obj.(type) {
	case *common.ObjString, *bytecode.ObjNative:
		// no outgoing references
	case *bytecode.Function:
		v.markObject(o.Name)
		if o.Chunk != nil && o.Chunk.Constants != nil {
			for _, constant := range o.Chunk.Constants.Values {
				v.markValue(constant)
			}
		}
	case *bytecode.Closure:
		v.markObject(o.Function)
		for i := 0; i < o.UpvalueCount; i++ {
			v.markObject(o.Upvalues[i])
		}
	case *bytecode.Upvalue:
		v.markValue(o.Closed)
	case *bytecode.Class:
		v.markObject(o.Name)
		v.markTable(&o.Methods)
	case *bytecode.Instance:
		v.markObject(o.Klass)
		v.markTable(&o.Fields)
	case *bytecode.BoundMethod:
		v.markValue(o.Receiver)
		v.markObject(o.Method)
	}
}

func (v *VM) sweep() {
	var previous common.Obj
	object := v.objects

	for object != nil {
		header := object.Header()
		if header.Marked {
			header.Marked = false
			previous = object
			object = header.Next
			continue
		}

		unreached := object
		object = header.Next
		if previous != nil {
			previous.Header().Next = object
		} else {
			v.objects = object
		}
		v.freeObject(unreached)
	}
}

func (v *VM) freeObject(obj common.Obj) {
	switch o := obj.(type) {
	case *common.ObjString:
		v.bytesAllocated -= sizeOfString(o)
	case *bytecode.Function:
		v.bytesAllocated -= sizeOfFunction(o)
	case *bytecode.Closure:
		v.bytesAllocated -= sizeOfClosure(o)
	case *bytecode.Upvalue:
		v.bytesAllocated -= sizeOfUpvalue()
	case *bytecode.ObjNative:
		v.bytesAllocated -= sizeOfNative()
	case *bytecode.Class:
		v.bytesAllocated -= sizeOfClass()
	case *bytecode.Instance:
		v.bytesAllocated -= sizeOfInstance()
	case *bytecode.BoundMethod:
		v.bytesAllocated -= sizeOfBoundMethod()
	}
	if v.bytesAllocated < 0 {
		v.bytesAllocated = 0
	}
}

func (v *VM) freeObjects() {
	object := v.objects
	for object != nil {
		next := object.Header().Next
		v.freeObject(object)
		object = next
	}
	v.objects = nil
	v.grayStack = nil
	v.bytesAllocated = 0
}
