package bytecode

import "golox/pkg/vm"

type Chunk struct {
	Code      []byte
	Constants *ValueArray
	Lines     []LineStart
}

func NewChunk() *Chunk {
	return &Chunk{
		Code:      make([]byte, 0, 8),
		Constants: NewValueArray(),
		Lines:     make([]LineStart, 0, 8),
	}
}

func (c *Chunk) Write(code byte, line int) {
	c.Code = append(c.Code, code)

	if len(c.Lines) > 0 && c.Lines[len(c.Lines)-1].Line == line {
		c.Lines[len(c.Lines)-1].Count++
	} else {
		c.Lines = append(c.Lines, LineStart{
			Line:  line,
			Count: 1,
		})
	}
}

func (c *Chunk) AddConstant(constant vm.Value) int {
	c.Constants.Write(constant)
	return c.Constants.Size() - 1
}

func (c *Chunk) Free() {
	c.Code = nil
	c.Constants = nil
	c.Lines = nil
}
