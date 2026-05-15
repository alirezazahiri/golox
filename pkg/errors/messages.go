package errors

const (
	UnexpectedCharacterErrorMessage string = "unexpected character"
	UnterminatedStringErrorMessage  string = "unterminated string"

	TooManyConstantsInOneChunk  string = "too many constants in one chunk"
	ExpectRParenAfterExpression string = "expected ')' after expression"
	ExpectExpression            string = "expect expression"
	OperandsMustBeNumbers       string = "operands must be numbers"
	OperandMustBeNumber         string = "operand must be a number"
)
