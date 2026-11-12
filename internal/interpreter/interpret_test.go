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

func TestClassInstanceFields(t *testing.T) {
	source := `
class Pair {}
var pair = Pair();
pair.first = 1;
pair.second = 2;
print pair.first;
print pair.second;
`
	result := compileAndRun(t, source)
	if result != vm.InterpretOk {
		t.Fatalf("expected InterpretOk, got %v", result)
	}
}

func TestUndefinedProperty(t *testing.T) {
	source := `
class Foo {}
var foo = Foo();
print foo.missing;
`
	result := compileAndRun(t, source)
	if result != vm.InterpretRuntimeError {
		t.Fatalf("expected InterpretRuntimeError, got %v", result)
	}
}

func TestMethodsAndInit(t *testing.T) {
	source := `
class Box {
  init(n) {
    this.n = n;
  }
  get() {
    return this.n;
  }
}
var b = Box(42);
print b.get();
`
	result := compileAndRun(t, source)
	if result != vm.InterpretOk {
		t.Fatalf("expected InterpretOk, got %v", result)
	}
}

func TestInitArity(t *testing.T) {
	source := `
class Box {
  init(n) {
    this.n = n;
  }
}
Box();
`
	result := compileAndRun(t, source)
	if result != vm.InterpretRuntimeError {
		t.Fatalf("expected InterpretRuntimeError for wrong init arity, got %v", result)
	}
}

func TestReturnValueFromInitializer(t *testing.T) {
	source := `
class Foo {
  init() {
    return 1;
  }
}
`
	result := compileAndRun(t, source)
	if result != vm.InterpretCompileError {
		t.Fatalf("expected InterpretCompileError, got %v", result)
	}
}

func TestThisOutsideClass(t *testing.T) {
	result := compileAndRun(t, `print this;`)
	if result != vm.InterpretCompileError {
		t.Fatalf("expected InterpretCompileError, got %v", result)
	}
}

func TestInheritanceAndSuper(t *testing.T) {
	source := `
class A {
  method() {
    print "A";
  }
}
class B < A {
  method() {
    super.method();
    print "B";
  }
}
B().method();
`
	result := compileAndRun(t, source)
	if result != vm.InterpretOk {
		t.Fatalf("expected InterpretOk, got %v", result)
	}
}

func TestInheritFromNonClass(t *testing.T) {
	source := `
var NotAClass = 123;
class Sub < NotAClass {}
`
	result := compileAndRun(t, source)
	if result != vm.InterpretRuntimeError {
		t.Fatalf("expected InterpretRuntimeError, got %v", result)
	}
}

func TestInheritFromSelf(t *testing.T) {
	source := `
class Oops < Oops {}
`
	result := compileAndRun(t, source)
	if result != vm.InterpretCompileError {
		t.Fatalf("expected InterpretCompileError, got %v", result)
	}
}

func TestSuperOutsideSubclass(t *testing.T) {
	source := `
class A {
  method() {
    super.method();
  }
}
`
	result := compileAndRun(t, source)
	if result != vm.InterpretCompileError {
		t.Fatalf("expected InterpretCompileError, got %v", result)
	}
}

func TestSuperInit(t *testing.T) {
	source := `
class A {
  init(n) {
    this.n = n;
  }
}
class B < A {
  init(n) {
    super.init(n);
  }
}
var b = B(7);
print b.n;
`
	result := compileAndRun(t, source)
	if result != vm.InterpretOk {
		t.Fatalf("expected InterpretOk, got %v", result)
	}
}
