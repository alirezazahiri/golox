package parser

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"golox/internal/chunk"
	"golox/internal/scanner"
	"golox/internal/vm"
	"golox/pkg/common"
	"golox/pkg/errors"
)

func ident(name string) scanner.Token {
	return scanner.NewToken(scanner.TOKEN_IDENTIFIER, 0, len(name), 1, name)
}

func testParser() *Parser {
	return New(scanner.New(""), chunk.New(), vm.New())
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe: %v", err)
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	return buf.String()
}

func compileSource(t *testing.T, source string) (*Parser, *chunk.Chunk, string) {
	t.Helper()

	var (
		p *Parser
		c *chunk.Chunk
	)

	out := captureStdout(t, func() {
		c = chunk.New()
		p = New(scanner.New(source), c, vm.New())
		p.Advance()
		for !p.Match(scanner.TOKEN_EOF) {
			p.Declaration()
		}
		p.EmitReturn()
	})

	return p, c, out
}

func assertCode(t *testing.T, got []byte, want ...byte) {
	t.Helper()
	if bytes.Equal(got, want) {
		return
	}
	t.Fatalf("bytecode mismatch\ngot  %v\nwant %v", got, want)
}

func assertClean(t *testing.T, p *Parser, out string) {
	t.Helper()
	if p.HadError || out != "" {
		t.Fatalf("unexpected compile error (hadError=%v)\n%s", p.HadError, out)
	}
	if p.compiler.ScopeDepth != 0 || p.compiler.LocalCount != 0 {
		t.Fatalf(
			"after compilation scope depth=%d local count=%d, want both 0",
			p.compiler.ScopeDepth,
			p.compiler.LocalCount,
		)
	}
}

func TestNewCompilerStartsAtGlobalScope(t *testing.T) {
	c := NewCompiler()

	if c.ScopeDepth != 0 {
		t.Fatalf("scope depth %d, want 0", c.ScopeDepth)
	}
	if c.LocalCount != 0 {
		t.Fatalf("local count %d, want 0", c.LocalCount)
	}
	for i, local := range c.Locals {
		if local != nil {
			t.Fatalf("slot %d is occupied before any declaration", i)
		}
	}
}

func TestAddLocalAndMarkInitialized(t *testing.T) {
	c := NewCompiler()
	c.ScopeDepth = 3

	c.addLocal(ident("a"))
	c.addLocal(ident("b"))

	if c.LocalCount != 2 {
		t.Fatalf("local count %d, want 2", c.LocalCount)
	}
	if c.Locals[0].Name.Lexeme != "a" || c.Locals[1].Name.Lexeme != "b" {
		t.Fatalf("locals stored as %q, %q", c.Locals[0].Name.Lexeme, c.Locals[1].Name.Lexeme)
	}
	if c.Locals[0].Depth != -1 || c.Locals[1].Depth != -1 {
		t.Fatalf(
			"new locals must be uninitialized, depths %d and %d",
			c.Locals[0].Depth,
			c.Locals[1].Depth,
		)
	}

	c.markInitialized()

	if c.Locals[0].Depth != -1 {
		t.Fatalf("markInitialized rewrote an older local to depth %d", c.Locals[0].Depth)
	}
	if c.Locals[1].Depth != 3 {
		t.Fatalf("newest local depth %d, want current scope 3", c.Locals[1].Depth)
	}

	c.ScopeDepth = 1
	if c.Locals[1].Depth != 3 {
		t.Fatalf("stored depth changed when the active scope changed")
	}

	c.markInitialized()
	if c.Locals[1].Depth != 1 {
		t.Fatalf("depth after the scope moved is %d, want 1", c.Locals[1].Depth)
	}
	if c.Locals[0].Depth != -1 {
		t.Fatalf("older local depth %d, want -1", c.Locals[0].Depth)
	}
}

func TestResolveLocalSelectsInnermostName(t *testing.T) {
	c := NewCompiler()
	c.ScopeDepth = 1
	c.addLocal(ident("a"))
	c.markInitialized()
	c.addLocal(ident("b"))
	c.markInitialized()

	c.ScopeDepth = 2
	c.addLocal(ident("a"))
	c.markInitialized()
	c.addLocal(ident("ab"))
	c.markInitialized()

	slot, ok := c.resolveLocal(ident("a"))
	if !ok || slot != 2 {
		t.Fatalf("resolve a = (%d, %v), want inner slot 2", slot, ok)
	}

	slot, ok = c.resolveLocal(ident("b"))
	if !ok || slot != 1 {
		t.Fatalf("resolve b = (%d, %v), want slot 1", slot, ok)
	}

	slot, ok = c.resolveLocal(ident("ab"))
	if !ok || slot != 3 {
		t.Fatalf("resolve ab = (%d, %v), want slot 3", slot, ok)
	}

	// Same spelling in a different token still matches. A prefix does not.
	elsewhere := scanner.NewToken(scanner.TOKEN_IDENTIFIER, 40, 1, 9, "a")
	slot, ok = c.resolveLocal(elsewhere)
	if !ok || slot != 2 {
		t.Fatalf("resolve a from another token = (%d, %v), want slot 2", slot, ok)
	}

	slot, ok = c.resolveLocal(ident("A"))
	if ok || slot != 0 {
		t.Fatalf("resolve A = (%d, %v), want missing (0, false)", slot, ok)
	}

	slot, ok = c.resolveLocal(ident("missing"))
	if ok || slot != 0 {
		t.Fatalf("resolve missing = (%d, %v), want (0, false)", slot, ok)
	}
}

func TestResolveLocalReportsUseInsideInitializer(t *testing.T) {
	c := NewCompiler()
	c.ScopeDepth = 1
	c.addLocal(ident("a"))
	c.markInitialized()

	c.ScopeDepth = 2
	c.addLocal(ident("a"))

	var (
		slot int
		ok   bool
	)
	out := captureStdout(t, func() {
		slot, ok = c.resolveLocal(ident("a"))
	})

	if !ok || slot != 1 {
		t.Fatalf("uninitialized lookup = (%d, %v), want the new local at slot 1", slot, ok)
	}
	if !strings.Contains(out, "can't read local variable in its own initializer") {
		t.Fatalf("stdout %q, want the own-initializer diagnostic", out)
	}

	out = captureStdout(t, func() {
		slot, ok = c.resolveLocal(ident("nope"))
	})
	if ok || slot != 0 {
		t.Fatalf("missing lookup = (%d, %v), want (0, false)", slot, ok)
	}
	if out != "" {
		t.Fatalf("missing name produced a diagnostic: %s", out)
	}

	c.markInitialized()
	out = captureStdout(t, func() {
		slot, ok = c.resolveLocal(ident("a"))
	})
	if !ok || slot != 1 {
		t.Fatalf("lookup after init = (%d, %v), want slot 1", slot, ok)
	}
	if out != "" {
		t.Fatalf("initialized local produced a diagnostic: %s", out)
	}
}

func TestAddLocalCapsTheLocalTable(t *testing.T) {
	c := NewCompiler()
	c.ScopeDepth = 1

	for i := 0; i < uInt8Count-1; i++ {
		c.addLocal(ident("v"))
	}
	c.addLocal(ident("last"))

	if c.LocalCount != uInt8Count {
		t.Fatalf("local count %d, want %d", c.LocalCount, uInt8Count)
	}

	slot, ok := c.resolveLocal(ident("last"))
	if !ok || slot != uInt8Count-1 {
		t.Fatalf("last legal slot = (%d, %v), want %d", slot, ok, uInt8Count-1)
	}

	out := captureStdout(t, func() {
		c.addLocal(ident("overflow"))
	})

	if c.LocalCount != uInt8Count {
		t.Fatalf("overflow grew the table to %d", c.LocalCount)
	}
	if c.Locals[uInt8Count-1].Name.Lexeme != "last" {
		t.Fatalf("overflow overwrote slot %d", uInt8Count-1)
	}
	if !strings.Contains(out, errors.TooManyLocalVariables) {
		t.Fatalf("stdout %q, want the too-many-locals diagnostic", out)
	}
}

func TestGlobalDeclarationsLeaveTheLocalTableEmpty(t *testing.T) {
	p, c, out := compileSource(t, "var a; var b = 1;")
	assertClean(t, p, out)
	assertCode(t, c.Code,
		byte(common.OpNil),
		byte(common.OpDefineGlobal), 0,
		byte(common.OpConstant), 2,
		byte(common.OpDefineGlobal), 1,
		byte(common.OpReturn),
	)
}

func TestLocalsInABlockGetSequentialSlots(t *testing.T) {
	p, c, out := compileSource(t, "{ var a = 1; var b = 2; var c = 3; print a; print b; print c; }")
	assertClean(t, p, out)
	assertCode(t, c.Code,
		byte(common.OpConstant), 0,
		byte(common.OpConstant), 1,
		byte(common.OpConstant), 2,
		byte(common.OpGetLocal), 0,
		byte(common.OpPrint),
		byte(common.OpGetLocal), 1,
		byte(common.OpPrint),
		byte(common.OpGetLocal), 2,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpPop),
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestOmittedInitializerStoresNil(t *testing.T) {
	p, c, out := compileSource(t, "{ var a; print a; }")
	assertClean(t, p, out)
	assertCode(t, c.Code,
		byte(common.OpNil),
		byte(common.OpGetLocal), 0,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestNestedBlocksPreserveOuterLocals(t *testing.T) {
	source := "{ var a = 1; { var b = 2; { print b; print a; } print a; } print a; }"
	p, c, out := compileSource(t, source)
	assertClean(t, p, out)
	assertCode(t, c.Code,
		byte(common.OpConstant), 0,
		byte(common.OpConstant), 1,
		byte(common.OpGetLocal), 1,
		byte(common.OpPrint),
		byte(common.OpGetLocal), 0,
		byte(common.OpPrint),
		byte(common.OpGetLocal), 0,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpGetLocal), 0,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestShadowedNamesResolveInnermostFirst(t *testing.T) {
	source := "{ var a = 1; { var a = 2; { var a = 3; print a; } print a; } print a; }"
	p, c, out := compileSource(t, source)
	assertClean(t, p, out)
	assertCode(t, c.Code,
		byte(common.OpConstant), 0,
		byte(common.OpConstant), 1,
		byte(common.OpConstant), 2,
		byte(common.OpGetLocal), 2,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpGetLocal), 1,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpGetLocal), 0,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestDuplicateLocalInOneScopeIsRejected(t *testing.T) {
	p, c, out := compileSource(t, "{ var a = 1; var a = 2; }")

	if !p.HadError {
		t.Fatal("expected a duplicate-name error")
	}
	if p.PanicMode {
		t.Fatal("parser stayed in panic mode after the declaration")
	}
	if p.compiler.ScopeDepth != 0 || p.compiler.LocalCount != 0 {
		t.Fatalf(
			"scope depth %d local count %d, want a restored global scope",
			p.compiler.ScopeDepth,
			p.compiler.LocalCount,
		)
	}
	if !strings.Contains(out, errors.AlreadyVariableWithThisNameInThisScope) {
		t.Fatalf("stdout %q, want the duplicate-name diagnostic", out)
	}

	// The declaration is still compiled, then both locals are popped with the block.
	assertCode(t, c.Code,
		byte(common.OpConstant), 0,
		byte(common.OpConstant), 1,
		byte(common.OpPop),
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestUninitializedLocalStillReservesItsName(t *testing.T) {
	p := testParser()
	p.compiler.ScopeDepth = 1
	p.Previous = ident("a")
	p.declareVariable()

	if p.compiler.LocalCount != 1 || p.compiler.Locals[0].Depth != -1 {
		t.Fatalf(
			"declared local count %d depth %d, want one uninitialized local",
			p.compiler.LocalCount,
			p.compiler.Locals[0].Depth,
		)
	}

	p.Previous = ident("a")
	out := captureStdout(t, func() {
		p.declareVariable()
	})

	if !p.HadError {
		t.Fatal("expected a duplicate while the first local is still uninitialized")
	}
	if !strings.Contains(out, errors.AlreadyVariableWithThisNameInThisScope) {
		t.Fatalf("stdout %q, want the duplicate-name diagnostic", out)
	}
	if p.compiler.LocalCount != 2 {
		t.Fatalf("local count %d, want 2 after the rejected declaration is still recorded", p.compiler.LocalCount)
	}
}

func TestNamesThatDifferByCaseAreDifferentLocals(t *testing.T) {
	p, c, out := compileSource(t, "{ var a = 1; var A = 2; print a; print A; }")
	assertClean(t, p, out)
	assertCode(t, c.Code,
		byte(common.OpConstant), 0,
		byte(common.OpConstant), 1,
		byte(common.OpGetLocal), 0,
		byte(common.OpPrint),
		byte(common.OpGetLocal), 1,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestSiblingScopesMayReuseANameAndSlot(t *testing.T) {
	source := "{ var a = 1; print a; } { var a = 2; print a; }"
	p, c, out := compileSource(t, source)
	assertClean(t, p, out)
	assertCode(t, c.Code,
		byte(common.OpConstant), 0,
		byte(common.OpGetLocal), 0,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpConstant), 1,
		byte(common.OpGetLocal), 0,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestLocalsDoNotEscapeTheirBlock(t *testing.T) {
	p, c, out := compileSource(t, "{ var a = 1; } print a;")
	assertClean(t, p, out)
	assertCode(t, c.Code,
		byte(common.OpConstant), 0,
		byte(common.OpPop),
		byte(common.OpGetGlobal), 1,
		byte(common.OpPrint),
		byte(common.OpReturn),
	)
}

func TestUseBeforeDeclarationResolvesAsGlobal(t *testing.T) {
	p, c, out := compileSource(t, "{ print a; var a = 1; print a; }")
	assertClean(t, p, out)
	assertCode(t, c.Code,
		byte(common.OpGetGlobal), 0,
		byte(common.OpPrint),
		byte(common.OpConstant), 1,
		byte(common.OpGetLocal), 0,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestInitializerSeesEarlierLocalsButNotItself(t *testing.T) {
	p, c, out := compileSource(t, "{ var a = 1; var b = a; print b; }")
	assertClean(t, p, out)
	assertCode(t, c.Code,
		byte(common.OpConstant), 0,
		byte(common.OpGetLocal), 0,
		byte(common.OpGetLocal), 1,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpPop),
		byte(common.OpReturn),
	)

	p, c, out = compileSource(t, "{ var a = a; }")
	if p.compiler.ScopeDepth != 0 || p.compiler.LocalCount != 0 {
		t.Fatalf(
			"scope depth %d local count %d, want a restored global scope",
			p.compiler.ScopeDepth,
			p.compiler.LocalCount,
		)
	}
	if !strings.Contains(out, "can't read local variable in its own initializer") {
		t.Fatalf("stdout %q, want the own-initializer diagnostic", out)
	}
	assertCode(t, c.Code,
		byte(common.OpGetLocal), 0,
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestInnerDeclarationHidesOuterInItsInitializer(t *testing.T) {
	source := "{ var a = 1; { var a = a; } }"
	p, c, out := compileSource(t, source)

	if p.compiler.ScopeDepth != 0 || p.compiler.LocalCount != 0 {
		t.Fatalf(
			"scope depth %d local count %d, want a restored global scope",
			p.compiler.ScopeDepth,
			p.compiler.LocalCount,
		)
	}
	if strings.Count(out, "can't read local variable in its own initializer") != 1 {
		t.Fatalf("stdout %q, want exactly one own-initializer diagnostic", out)
	}
	assertCode(t, c.Code,
		byte(common.OpConstant), 0,
		byte(common.OpGetLocal), 1,
		byte(common.OpPop),
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestLocalShadowsGlobalReadsAndWrites(t *testing.T) {
	source := "var a = 1; { var a = 2; a = 3; } a = 4;"
	p, c, out := compileSource(t, source)
	assertClean(t, p, out)
	assertCode(t, c.Code,
		byte(common.OpConstant), 1,
		byte(common.OpDefineGlobal), 0,
		byte(common.OpConstant), 2,
		byte(common.OpConstant), 3,
		byte(common.OpSetLocal), 0,
		byte(common.OpPop),
		byte(common.OpPop),
		byte(common.OpConstant), 5,
		byte(common.OpSetGlobal), 4,
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestInnerAssignmentUpdatesTheOuterLocal(t *testing.T) {
	p, c, out := compileSource(t, "{ var a = 1; { a = 2; } print a; }")
	assertClean(t, p, out)
	assertCode(t, c.Code,
		byte(common.OpConstant), 0,
		byte(common.OpConstant), 1,
		byte(common.OpSetLocal), 0,
		byte(common.OpPop),
		byte(common.OpGetLocal), 0,
		byte(common.OpPrint),
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestEndScopeDropsLocalsFromTheScopeThatEnded(t *testing.T) {
	p := testParser()
	p.compiler.ScopeDepth = 2
	p.compiler.addLocal(ident("outer"))
	p.compiler.markInitialized()

	p.compiler.ScopeDepth = 3
	p.compiler.addLocal(ident("left"))
	p.compiler.markInitialized()
	p.compiler.addLocal(ident("right"))
	p.compiler.markInitialized()

	p.endScope()

	if p.compiler.ScopeDepth != 2 {
		t.Fatalf("scope depth %d, want 2", p.compiler.ScopeDepth)
	}
	if p.compiler.LocalCount != 1 || p.compiler.Locals[0].Name.Lexeme != "outer" {
		t.Fatalf("remaining locals %d, want only outer", p.compiler.LocalCount)
	}
	if p.compiler.Locals[0].Depth != 2 {
		t.Fatalf("outer depth %d, want 2", p.compiler.Locals[0].Depth)
	}
	assertCode(t, p.chunk.Code, byte(common.OpPop), byte(common.OpPop))

	p.endScope()

	if p.compiler.ScopeDepth != 1 || p.compiler.LocalCount != 0 {
		t.Fatalf(
			"after the outer scope, depth %d count %d, want depth 1 and no locals",
			p.compiler.ScopeDepth,
			p.compiler.LocalCount,
		)
	}
	assertCode(t, p.chunk.Code, byte(common.OpPop), byte(common.OpPop), byte(common.OpPop))

	p = testParser()
	p.compiler.ScopeDepth = 1
	p.compiler.addLocal(ident("pending"))
	p.endScope()

	if p.compiler.ScopeDepth != 0 {
		t.Fatalf("scope depth %d, want 0", p.compiler.ScopeDepth)
	}
	if p.compiler.LocalCount != 1 || p.compiler.Locals[0].Depth != -1 {
		t.Fatalf(
			"uninitialized local count %d depth %d, want it to stay uninitialized",
			p.compiler.LocalCount,
			p.compiler.Locals[0].Depth,
		)
	}
	if len(p.chunk.Code) != 0 {
		t.Fatalf("uninitialized local emitted %v", p.chunk.Code)
	}
}

func TestMissingClosingBraceStillLeavesGlobalScope(t *testing.T) {
	p, c, out := compileSource(t, "{ var a = 1;")

	if !p.HadError {
		t.Fatal("expected an error for the missing '}'")
	}
	if p.PanicMode {
		t.Fatal("parser stayed in panic mode")
	}
	if p.compiler.ScopeDepth != 0 || p.compiler.LocalCount != 0 {
		t.Fatalf(
			"scope depth %d local count %d, want a restored global scope",
			p.compiler.ScopeDepth,
			p.compiler.LocalCount,
		)
	}
	if !strings.Contains(out, errors.ExpectRightBraceAfterBlock) {
		t.Fatalf("stdout %q, want the missing-'}' diagnostic", out)
	}
	assertCode(t, c.Code,
		byte(common.OpConstant), 0,
		byte(common.OpPop),
		byte(common.OpReturn),
	)
}

func TestEmptyBlocksDoNotLeaveScopeBehind(t *testing.T) {
	p, c, out := compileSource(t, "{ { { } } }")
	assertClean(t, p, out)
	assertCode(t, c.Code, byte(common.OpReturn))
}
