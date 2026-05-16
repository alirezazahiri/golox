package table

import "golox/pkg/common"

func (t *Table) FindString(chars string, hash uint32) *common.ObjString {
	if t.Count == 0 {
		return nil
	}

	index := int(hash % uint32(t.Capacity))

	for {
		entry := &t.Entries[index]

		if entry.Key == nil {
			if entry.Value.IsNil() {
				return nil
			}
		} else if entry.Key.Hash == hash &&
			entry.Key.Content == chars {
			return entry.Key
		}

		index = (index + 1) % t.Capacity
	}
}
