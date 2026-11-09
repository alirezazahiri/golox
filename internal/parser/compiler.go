package parser

import (
	"fmt"

	"golox/internal/scanner"
	"golox/pkg/bytecode"
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
	Enclosing  *Compiler
	Function   *bytecode.Function
	Type       FunctionType
	Locals     [uInt8Count]*Local
	LocalCount int
	ScopeDepth int
}

func (p *Parser) beginFunctionCompiler(ty FunctionType) {
	enclosing := p.Compiler
	compiler := &Compiler{
		Enclosing:  enclosing,
		LocalCount: 0,
		ScopeDepth: 0,
		Function:   bytecode.NewFunction(),
		Type:       ty,
	}
	if ty != TypeScript {
		compiler.Function.Name = p.vm.InternString(p.Previous.Lexeme)
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
		LocalCount: 0,
		ScopeDepth: 0,
		Function:   bytecode.NewFunction(),
		Type:       ty,
	}
	current.Locals[current.LocalCount] = &Local{Depth: 0}
	current.LocalCount++

	return current
}

func (c *Compiler) addLocal(name scanner.Token) {
	if c.LocalCount >= uInt8Count {
		c.error(errors.TooManyLocalVariables)
		return
	}
	c.Locals[c.LocalCount] = &Local{
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

func (c *Compiler) resolveLocal(name scanner.Token) (int, bool) {
	for i := c.LocalCount - 1; i >= 0; i-- {
		local := c.Locals[i]
		if identifiersEqual(&local.Name, &name) {
			if local.Depth == -1 {
				c.error(errors.CannotAccessVariableInItsInitializer)
			}
			return i, true
		}
	}
	return 0, false
}

func (c *Compiler) error(msg string) {
	fmt.Printf("[line %d] Error: %s\n", c.ScopeDepth, msg)
}
