package parser

import (
	"golox/pkg/bytecode"
	"golox/pkg/common"
)

// Allocator provides heap allocation for the compiler (strings + functions).
type Allocator interface {
	InternString(chars string) *common.ObjString
	AllocateFunction() *bytecode.Function
}
