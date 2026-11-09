package interpreter

import (
	"testing"

	"golox/internal/vm"
)

func compileAndRun(t *testing.T, source string) vm.InterpretResult {
	t.Helper()
	v := vm.New()
	t.Cleanup(v.Free)
	return Interpret(v, source)
}

func TestArithmetic(t *testing.T) {
	result := compileAndRun(t, "print 1 + 2 * 3;")
	if result != vm.InterpretOk {
		t.Fatalf("expected InterpretOk, got %v", result)
	}
}

func TestLocals(t *testing.T) {
	source := `
{
  var a = 10;
  var b = 20;
  print a + b;
}
`
	result := compileAndRun(t, source)
	if result != vm.InterpretOk {
		t.Fatalf("expected InterpretOk, got %v", result)
	}
}

func TestIfWhile(t *testing.T) {
	source := `
var i = 0;
while (i < 3) {
  i = i + 1;
}
if (i == 3) {
  print i;
}
`
	result := compileAndRun(t, source)
	if result != vm.InterpretOk {
		t.Fatalf("expected InterpretOk, got %v", result)
	}
}

func TestSingleParamFunctionDeclaration(t *testing.T) {
	source := `
fun identity(x) {
  print x;
}
print identity;
`
	result := compileAndRun(t, source)
	if result != vm.InterpretOk {
		t.Fatalf("expected InterpretOk, got %v", result)
	}
}

func TestMultiParamFunctionDeclaration(t *testing.T) {
	source := `
fun add(a, b) {
  print a;
}
print add;
`
	result := compileAndRun(t, source)
	if result != vm.InterpretOk {
		t.Fatalf("expected InterpretOk, got %v", result)
	}
}

func TestManyGlobalsSurviveResize(t *testing.T) {
	source := `
var g0 = 0;
var g1 = 1;
var g2 = 2;
var g3 = 3;
var g4 = 4;
var g5 = 5;
var g6 = 6;
var g7 = 7;
var g8 = 8;
var g9 = 9;
print g0;
print g9;
`
	result := compileAndRun(t, source)
	if result != vm.InterpretOk {
		t.Fatalf("expected InterpretOk after global table growth, got %v", result)
	}
}

func TestStringConcat(t *testing.T) {
	result := compileAndRun(t, `print "hel" + "lo";`)
	if result != vm.InterpretOk {
		t.Fatalf("expected InterpretOk, got %v", result)
	}
}

func TestCompileError(t *testing.T) {
	result := compileAndRun(t, "var;")
	if result != vm.InterpretCompileError {
		t.Fatalf("expected InterpretCompileError, got %v", result)
	}
}
