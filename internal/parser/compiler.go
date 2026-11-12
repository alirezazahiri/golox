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
	TypeMethod
	TypeInitializer
)

const uInt8Count int = 256

type Local struct {
	Name       scanner.Token
	Depth      int
	IsCaptured bool
}

type Upvalue struct {
	Index   uint8
	IsLocal bool
}

type ClassCompiler struct {
	Enclosing     *ClassCompiler
	HasSuperclass bool
}

type Compiler struct {
	Enclosing     *Compiler
	Function      *bytecode.Function
	Type          FunctionType
	Locals        [uInt8Count]Local
	LocalCount    int
	Upvalues      [uInt8Count]Upvalue
	UpvalueCount  int
	ScopeDepth    int
	constantIndex map[*common.ObjString]uint8
}

func (p *Parser) beginFunctionCompiler(ty FunctionType) {
	enclosing := p.Compiler
	compiler := &Compiler{
		Enclosing:     enclosing,
		LocalCount:    0,
		ScopeDepth:    0,
		Function:      p.allocator.AllocateFunction(),
		Type:          ty,
		constantIndex: make(map[*common.ObjString]uint8),
	}
	if ty != TypeScript {
		compiler.Function.Name = p.allocator.InternString(p.Previous.Lexeme)
	}
	if ty == TypeInitializer {
		compiler.Function.IsInitializer = true
	}

	local := Local{Depth: 0}
	if ty != TypeFunction && ty != TypeScript {
		local.Name = syntheticToken("this")
	}
	compiler.Locals[compiler.LocalCount] = local
	compiler.LocalCount++

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

func NewCompilerWithFunction(ty FunctionType, function *bytecode.Function) *Compiler {
	compiler := &Compiler{
		LocalCount:    0,
		ScopeDepth:    0,
		Function:      function,
		Type:          ty,
		constantIndex: make(map[*common.ObjString]uint8),
	}
	compiler.Locals[compiler.LocalCount] = Local{Depth: 0}
	compiler.LocalCount++
	return compiler
}

func (p *Parser) addLocal(name scanner.Token) {
	c := p.Compiler
	if c.LocalCount >= uInt8Count {
		p.error(errors.TooManyLocalVariables)
		return
	}
	c.Locals[c.LocalCount] = Local{
		Name:       name,
		Depth:      -1,
		IsCaptured: false,
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
	return p.resolveLocalIn(p.Compiler, name)
}

func syntheticToken(text string) scanner.Token {
	return scanner.Token{
		Type:   scanner.TOKEN_IDENTIFIER,
		Lexeme: text,
	}
}
