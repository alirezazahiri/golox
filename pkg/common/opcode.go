package common

type OpCode byte

const (
	OpReturn OpCode = iota
	OpConstant
	OpConstantLong
	OpEqual
	OpGreater
	OpLess
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpNegate
	OpPop
	OpPrint
	OpBang
	OpFalse
	OpTrue 
	OpNil
	OpDefineGlobal
	OpGetGlobal
	OpGetLocal
	OpSetGlobal
	OpSetLocal
)