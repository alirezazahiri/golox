package table

import "golox/pkg/common"

const tableMaxLoad = 0.75

func (t *Table) adjustCapacity(capacity int) {
	entries := make([]Entry, capacity)
	for i := range capacity {
		entries[i].Key = nil
		entries[i].Value = common.NilValue()
	}

	t.Count = 0
	for i := range t.Capacity {
		entry := &t.Entries[i]
		if entry.Key == nil {
			continue
		}

		dest := t.findEntry(entry.Key)
		dest.Key = entry.Key
		dest.Value = entry.Value
		t.Count++
	}

	t.Entries = entries
	t.Capacity = capacity
}

func growCapacity(cap int) int {
	if cap < 8 {
		return 8
	}
	return cap * 2
}
