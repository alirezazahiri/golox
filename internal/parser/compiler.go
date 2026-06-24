package parser

import (
	"fmt"
	"golox/internal/scanner"
	"golox/pkg/errors"
)

const uInt8Count int = 256

type Local struct {
	Name  scanner.Token
	Depth int
}

type Compiler struct {
	Locals     [uInt8Count]*Local
	LocalCount int
	ScopeDepth int
}

func NewCompiler() *Compiler {
	return &Compiler{
		LocalCount: 0,
		ScopeDepth: 0,
	}
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
	c.Locals[c.LocalCount-1].Depth = c.ScopeDepth
}

func (c *Compiler) resolveLocal(name scanner.Token) (int, bool) {
	for i := c.LocalCount - 1; i >= 0; i-- {
		local := c.Locals[i]
		if local.Name.Lexeme == name.Lexeme {
			if local.Depth == -1 {
				// error: can't read local before initializer
			}
			return i, true
		}
	}
	return 0, false
}

func (c *Compiler) error(msg string) {
	fmt.Printf("[line %d] Error: %s\n", c.ScopeDepth, msg)
}
