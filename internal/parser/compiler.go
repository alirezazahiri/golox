package parser

import (
	"golox/internal/scanner"
	"golox/pkg/bytecode"
	"golox/pkg/common"
	"golox/pkg/errors"
)

type FunctionType = byte

const (
	TypeFunction FunctionType = iota + 1
	TypeScript
)

const uInt8Count int = 256

type Local struct {
	Name  scanner.Token
	Depth int
}

type Compiler struct {
	Enclosing     *Compiler
	Function      *bytecode.Function
	Type          FunctionType
	Locals        [uInt8Count]Local
	LocalCount    int
	ScopeDepth    int
	constantIndex map[*common.ObjString]uint8
}

func (p *Parser) beginFunctionCompiler(ty FunctionType) {
	enclosing := p.Compiler
	compiler := &Compiler{
		Enclosing:     enclosing,
		LocalCount:    0,
		ScopeDepth:    0,
		Function:      bytecode.NewFunction(),
		Type:          ty,
		constantIndex: make(map[*common.ObjString]uint8),
	}
	if ty != TypeScript {
		compiler.Function.Name = p.interner.InternString(p.Previous.Lexeme)
	}
	p.Compiler = compiler
	p.chunk = compiler.Function.Chunk
}

func (p *Parser) EndCompiler() *bytecode.Function {
	p.EmitReturn()
	function := p.Compiler.Function
	p.Compiler = p.Compiler.Enclosing
	p.chunk = p.Compiler.Function.Chunk
	return function
}

func NewCompiler(ty FunctionType) *Compiler {
	current := &Compiler{
		LocalCount:    0,
		ScopeDepth:    0,
		Function:      bytecode.NewFunction(),
		Type:          ty,
		constantIndex: make(map[*common.ObjString]uint8),
	}
	current.Locals[current.LocalCount] = Local{Depth: 0}
	current.LocalCount++

	return current
}

func (p *Parser) addLocal(name scanner.Token) {
	c := p.Compiler
	if c.LocalCount >= uInt8Count {
		p.error(errors.TooManyLocalVariables)
		return
	}
	c.Locals[c.LocalCount] = Local{
		Name:  name,
		Depth: -1,
	}
	c.LocalCount++
}

func (c *Compiler) markInitialized() {
	if c.ScopeDepth == 0 {
		return
	}
	c.Locals[c.LocalCount-1].Depth = c.ScopeDepth
}

func (p *Parser) resolveLocal(name scanner.Token) (int, bool) {
	c := p.Compiler
	for i := c.LocalCount - 1; i >= 0; i-- {
		local := &c.Locals[i]
		if identifiersEqual(&local.Name, &name) {
			if local.Depth == -1 {
				p.error(errors.CannotAccessVariableInItsInitializer)
			}
			return i, true
		}
	}
	return 0, false
}
