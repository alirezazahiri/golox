package chunk

import "golox/pkg/common"

type Chunk struct {
	Code      []byte	
	Constants *common.ValueArray
	Lines     []LineStart
}

func New() *Chunk {
	return &Chunk{
		Code:      make([]byte, 0, 8),
		Constants: common.NewValueArray(),
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

func (c *Chunk) AddConstant(constant common.Value) int {
	c.Constants.Write(constant)
	return c.Constants.Size() - 1
}

func (c *Chunk) WriteConstant(value common.Value, line int) {
	index := c.AddConstant(value)

	if index <= 0xff {
		c.Write(byte(common.OpConstant), line)
		c.Write(byte(index), line)
	} else {
		c.Write(byte(common.OpConstantLong), line)

		c.Write(byte((index >> 16) & 0xff), line)
		c.Write(byte((index >> 8) & 0xff), line)
		c.Write(byte(index & 0xff), line)
	}
}


func (c *Chunk) Free() {
	c.Code = nil
	c.Constants = nil
	c.Lines = nil
}
