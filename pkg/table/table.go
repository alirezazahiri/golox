package table

import "golox/pkg/common"

type Entry struct {
	Key   *common.ObjString
	Value common.Value
}

type Table struct {
	Count    int
	Capacity int
	Entries  []Entry
}

func New() *Table {
	return &Table{}
}

func (t *Table) Set(key *common.ObjString, value common.Value) bool {
	if float64(t.Count+1) > float64(t.Capacity)*tableMaxLoad {
		capacity := growCapacity(t.Capacity)
		t.adjustCapacity(capacity)
	}

	entry := t.findEntry(key)

	isNewKey := entry.Key == nil
	if isNewKey && entry.Value.IsNil() {
		t.Count++
	}

	entry.Key = key
	entry.Value = value

	return isNewKey
}

func (t *Table) Get(key *common.ObjString) (value common.Value, found bool) {
	if t.Count == 0 {
		return
	}

	e := t.findEntry(key)

	if e.Key == nil {
		return
	}

	return e.Value, true
}

func (t *Table) Delete(key *common.ObjString) bool {
	if t.Count == 0 {
		return false
	}

	e := t.findEntry(key)
	if e.Key == nil {
		return false
	}

	e.Key = nil
	e.Value = common.BoolValue(true)

	return true
}

// AddAll copies all entries from this table into to (used later for class inheritance).
func (from *Table) AddAll(to *Table) {
	for _, e := range from.Entries {
		if e.Key != nil {
			to.Set(e.Key, e.Value)
		}
	}
}

// RemoveWhite deletes entries whose keys were not marked by the GC.
func (t *Table) RemoveWhite() {
	for i := range t.Entries {
		e := &t.Entries[i]
		if e.Key != nil && !e.Key.Marked {
			t.Delete(e.Key)
		}
	}
}

func (t *Table) findEntry(key *common.ObjString) *Entry {
	return findEntryIn(t.Entries, t.Capacity, key)
}

func findEntryIn(entries []Entry, capacity int, key *common.ObjString) *Entry {
	index := int(key.Hash % uint32(capacity))
	var tombstone *Entry

	for {
		entry := &entries[index]

		if entry.Key == nil {
			if entry.Value.IsNil() {
				if tombstone != nil {
					return tombstone
				}
				return entry
			} else {
				if tombstone == nil {
					tombstone = entry
				}
			}
		} else if entry.Key == key {
			return entry
		}

		index = (index + 1) % capacity
	}
}
