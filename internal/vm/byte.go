package vm

func (v *VM) ReadByte() (byte, error) {
	frame := v.currentFrame()
	code := frame.Closure.Function.Chunk.Code
	b := code[frame.IP]
	frame.IP++
	return b, nil
}
