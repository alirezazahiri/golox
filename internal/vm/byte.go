package vm

func (v *VM) ReadByte() (byte, error) {
	b := v.Chunk.Code[v.IP]
	v.IP++
	return b, nil
}