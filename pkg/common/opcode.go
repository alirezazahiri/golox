package common

type OpCode byte

const (
	OpReturn OpCode = iota
	OpConstant
	OpConstantLong
)