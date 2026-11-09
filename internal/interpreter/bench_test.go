package interpreter

import (
	"testing"

	"golox/internal/vm"
)

const benchSource = `
var i = 0;
while (i < 100) {
  i = i + 1;
}
print i;
`

func BenchmarkCompile(b *testing.B) {
	v := vm.New()
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if compile(benchSource, v) == nil {
			b.Fatal("compile failed")
		}
	}
}

func BenchmarkRun(b *testing.B) {
	v := vm.New()
	function := compile(benchSource, v)
	if function == nil {
		b.Fatal("compile failed")
	}
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if result := v.InterpretFunction(function); result != vm.InterpretOk {
			b.Fatalf("run failed: %v", result)
		}
	}
}
