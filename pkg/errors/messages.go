package errors

const (
	UnexpectedCharacterErrorMessage string = "unexpected character"
	UnterminatedStringErrorMessage  string = "unterminated string"

	TooManyConstantsInOneChunk              string = "too many constants in one chunk"
	ExpectRParenAfterExpression             string = "expected ')' after expression"
	ExpectExpression                        string = "expected expression"
	ExpectSemicolonAfterPrintStatement      string = "expected ';' after print statement"
	ExpectSemicolonAfterExpressionStatement string = "expected ';' after expression statement"
	ExpectIdentifierAfterVar                string = "expected identifier after 'var'"
	ExpectEqualAfterVar                     string = "expected '=' after 'var'"
	ExpectSemicolonAfterVarDeclaration      string = "expected ';' after variable declaration"
	OperandsMustBeNumbers                   string = "operands must be numbers"
	OperandMustBeNumber                     string = "operand must be a number"
	ExpectedOperandToBeOfType               string = "expected operand to be of type %s, but got %s"
	ExpectedOperandsToBeOfType              string = "expected operands to be of type %s, but got (%s, %s)"
	FailedToSetValueToVariable              string = "failed to set value for the variable %s"
	UndefinedVariable                       string = "variable with name \"%s\" is not defined"
	
)
