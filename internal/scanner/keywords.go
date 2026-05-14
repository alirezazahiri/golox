package scanner

var keywords = map[string]TokenType{
	"and":    TOKEN_AND,
	"class":  TOKEN_CLASS,
	"else":   TOKEN_ELSE,
	"false":  TOKEN_FALSE,
	"fun":    TOKEN_FUN,
	"for":    TOKEN_FOR,
	"if":     TOKEN_IF,
	"nil":    TOKEN_NIL,
	"or":     TOKEN_OR,
	"print":  TOKEN_PRINT,
	"return": TOKEN_RETURN,
	"super":  TOKEN_SUPER,
	"this":   TOKEN_THIS,
	"true":   TOKEN_TRUE,
	"var":    TOKEN_VAR,
	"while":  TOKEN_WHILE,
	"set":    TOKEN_SET,
}