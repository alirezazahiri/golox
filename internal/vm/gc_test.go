package vm_test

import (
	"fmt"
	"strings"
	"testing"

	"golox/internal/interpreter"
	"golox/internal/vm"
	"golox/pkg/common"
)

func TestGCReclaimsInternedStrings(t *testing.T) {
	v := vm.New()
	t.Cleanup(v.Free)

	before := v.ObjectCount()

	const n = 50
	keys := make([]*common.ObjString, n)
	for i := 0; i < n; i++ {
		keys[i] = v.InternString(fmt.Sprintf("tmp-%d", i))
		v.Globals.Set(keys[i], common.NumberValue(float64(i)))
	}

	mid := v.ObjectCount()
	if mid < before+n {
		t.Fatalf("expected at least %d objects after interning, got %d", before+n, mid)
	}

	for _, k := range keys {
		v.Globals.Delete(k)
	}

	v.CollectGarbage()

	after := v.ObjectCount()
	reclaimed := mid - after
	if reclaimed < n {
		t.Fatalf("expected to reclaim at least %d objects, reclaimed %d (mid=%d after=%d)", n, reclaimed, mid, after)
	}
	if after < before {
		t.Fatalf("GC reclaimed natives; before=%d after=%d", before, after)
	}
}

func TestClosureSurvivesGC(t *testing.T) {
	v := vm.New()
	t.Cleanup(v.Free)

	source := `
fun make() {
  var x = 1;
  fun inner() {
    return x;
  }
  return inner;
}

var f = make();
`
	result := interpreter.Interpret(v, source)
	if result != vm.InterpretOk {
		t.Fatalf("setup failed: %v", result)
	}

	v.CollectGarbage()

	result = interpreter.Interpret(v, `print f();`)
	if result != vm.InterpretOk {
		t.Fatalf("closure call after GC failed: %v", result)
	}
}

func TestGCCollectionUnderAllocationPressure(t *testing.T) {
	v := vm.New()
	t.Cleanup(v.Free)
	v.SetNextGC(64)

	var b strings.Builder
	b.WriteString("var s = \"\";\n")
	b.WriteString("var i = 0;\n")
	b.WriteString("while (i < 200) {\n")
	b.WriteString("  s = s + \"x\";\n")
	b.WriteString("  i = i + 1;\n")
	b.WriteString("}\n")

	result := interpreter.Interpret(v, b.String())
	if result != vm.InterpretOk {
		t.Fatalf("expected InterpretOk under GC pressure, got %v", result)
	}
}
