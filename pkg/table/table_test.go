package table

import (
	"fmt"
	"testing"

	"golox/pkg/common"
)

func makeKey(chars string) *common.ObjString {
	return &common.ObjString{
		Content: chars,
		Length:  len(chars),
		Hash:    common.HashString(chars),
	}
}

func TestSetAndGet(t *testing.T) {
	tbl := New()
	key := makeKey("answer")

	isNew := tbl.Set(key, common.NumberValue(42))
	if !isNew {
		t.Fatal("expected first Set to report a new key")
	}

	value, found := tbl.Get(key)
	if !found {
		t.Fatal("expected key to be found after Set")
	}
	if !value.IsNumber() || value.As.Number != 42 {
		t.Fatalf("expected 42, got %v", value)
	}
}

func TestSetExistingKey(t *testing.T) {
	tbl := New()
	key := makeKey("x")

	tbl.Set(key, common.NumberValue(1))
	isNew := tbl.Set(key, common.NumberValue(2))
	if isNew {
		t.Fatal("expected second Set to report an existing key")
	}

	value, found := tbl.Get(key)
	if !found || value.As.Number != 2 {
		t.Fatalf("expected updated value 2, got found=%v value=%v", found, value)
	}
}

func TestDeleteAndTombstone(t *testing.T) {
	tbl := New()
	key := makeKey("gone")

	tbl.Set(key, common.NumberValue(7))
	if !tbl.Delete(key) {
		t.Fatal("expected Delete to succeed")
	}

	if _, found := tbl.Get(key); found {
		t.Fatal("expected key to be missing after Delete")
	}

	// Re-insert into a tombstone slot should still work.
	if !tbl.Set(key, common.NumberValue(8)) {
		t.Fatal("expected re-insert after delete to be a new key")
	}
	value, found := tbl.Get(key)
	if !found || value.As.Number != 8 {
		t.Fatalf("expected 8 after re-insert, got found=%v value=%v", found, value)
	}
}

func TestResizePreservesEntries(t *testing.T) {
	tbl := New()
	const n = 32

	keys := make([]*common.ObjString, n)
	for i := range n {
		keys[i] = makeKey(fmt.Sprintf("g%d", i))
		tbl.Set(keys[i], common.NumberValue(float64(i)))
	}

	if tbl.Capacity < 8 {
		t.Fatalf("expected capacity to grow, got %d", tbl.Capacity)
	}

	for i, key := range keys {
		value, found := tbl.Get(key)
		if !found {
			t.Fatalf("missing key %s after resize", key.Content)
		}
		if value.As.Number != float64(i) {
			t.Fatalf("key %s: expected %d, got %v", key.Content, i, value.As.Number)
		}
	}
}

func TestFindString(t *testing.T) {
	tbl := New()
	key := makeKey("hello")
	tbl.Set(key, common.NilValue())

	found := tbl.FindString("hello", key.Hash)
	if found != key {
		t.Fatal("expected FindString to return the interned key")
	}

	if tbl.FindString("missing", common.HashString("missing")) != nil {
		t.Fatal("expected FindString to miss unknown string")
	}
}
