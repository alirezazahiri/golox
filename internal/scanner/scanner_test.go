package scanner

import "testing"

func collectTokens(s *Scanner) []Token {
	var tokens []Token

	for {
		tok := s.ScanToken()
		tokens = append(tokens, tok)

		if tok.Type == TOKEN_EOF || tok.Type == TOKEN_ERROR {
			break
		}
	}

	return tokens
}

func tokenTypes(tokens []Token) []TokenType {
	types := make([]TokenType, len(tokens))

	for i, tok := range tokens {
		types[i] = tok.Type
	}

	return types
}

func assertTokenTypes(t *testing.T, got []Token, want []TokenType) {
	t.Helper()

	gotTypes := tokenTypes(got)

	if len(gotTypes) != len(want) {
		t.Fatalf("wrong token count\ngot  %v\nwant %v", gotTypes, want)
	}

	for i := range want {
		if gotTypes[i] != want[i] {
			t.Fatalf(
				"wrong token at index %d\ngot  %v\nwant %v",
				i,
				gotTypes[i],
				want[i],
			)
		}
	}
}

func TestSingleCharacterTokens(t *testing.T) {
	s := New("(){}.,-+;/*")

	tokens := collectTokens(s)

	assertTokenTypes(t, tokens, []TokenType{
		TOKEN_LEFT_PAREN,
		TOKEN_RIGHT_PAREN,
		TOKEN_LEFT_BRACE,
		TOKEN_RIGHT_BRACE,
		TOKEN_DOT,
		TOKEN_COMMA,
		TOKEN_MINUS,
		TOKEN_PLUS,
		TOKEN_SEMICOLON,
		TOKEN_SLASH,
		TOKEN_STAR,
		TOKEN_EOF,
	})
}

func TestOneOrTwoCharacterTokens(t *testing.T) {
	s := New("! != = == < <= > >=")

	tokens := collectTokens(s)

	assertTokenTypes(t, tokens, []TokenType{
		TOKEN_BANG,
		TOKEN_BANG_EQUAL,
		TOKEN_EQUAL,
		TOKEN_EQUAL_EQUAL,
		TOKEN_LESS,
		TOKEN_LESS_EQUAL,
		TOKEN_GREATER,
		TOKEN_GREATER_EQUAL,
		TOKEN_EOF,
	})
}

func TestIdentifiers(t *testing.T) {
	s := New("hello world foo_bar abc123")

	tokens := collectTokens(s)

	assertTokenTypes(t, tokens, []TokenType{
		TOKEN_IDENTIFIER,
		TOKEN_IDENTIFIER,
		TOKEN_IDENTIFIER,
		TOKEN_IDENTIFIER,
		TOKEN_EOF,
	})
}

func TestKeywords(t *testing.T) {
	s := New(`
and class else false
for fun if nil or
print return super this
true var while set
`)

	tokens := collectTokens(s)

	assertTokenTypes(t, tokens, []TokenType{
		TOKEN_AND,
		TOKEN_CLASS,
		TOKEN_ELSE,
		TOKEN_FALSE,
		TOKEN_FOR,
		TOKEN_FUN,
		TOKEN_IF,
		TOKEN_NIL,
		TOKEN_OR,
		TOKEN_PRINT,
		TOKEN_RETURN,
		TOKEN_SUPER,
		TOKEN_THIS,
		TOKEN_TRUE,
		TOKEN_VAR,
		TOKEN_WHILE,
		TOKEN_SET,
		TOKEN_EOF,
	})
}

func TestNumbers(t *testing.T) {
	s := New("123 456 789")

	tokens := collectTokens(s)

	assertTokenTypes(t, tokens, []TokenType{
		TOKEN_NUMBER,
		TOKEN_NUMBER,
		TOKEN_NUMBER,
		TOKEN_EOF,
	})
}

func TestStrings(t *testing.T) {
	s := New(`"hello" "world"`)

	tokens := collectTokens(s)

	assertTokenTypes(t, tokens, []TokenType{
		TOKEN_STRING,
		TOKEN_STRING,
		TOKEN_EOF,
	})
}

func TestComments(t *testing.T) {
	s := New(`
var a = 1; // comment
print a;
`)

	tokens := collectTokens(s)

	assertTokenTypes(t, tokens, []TokenType{
		TOKEN_VAR,
		TOKEN_IDENTIFIER,
		TOKEN_EQUAL,
		TOKEN_NUMBER,
		TOKEN_SEMICOLON,
		TOKEN_PRINT,
		TOKEN_IDENTIFIER,
		TOKEN_SEMICOLON,
		TOKEN_EOF,
	})
}

func TestWhitespace(t *testing.T) {
	s := New("   \r\t\n   var   ")

	tokens := collectTokens(s)

	assertTokenTypes(t, tokens, []TokenType{
		TOKEN_VAR,
		TOKEN_EOF,
	})
}

func TestUnexpectedCharacter(t *testing.T) {
	s := New("@")

	tokens := collectTokens(s)

	assertTokenTypes(t, tokens, []TokenType{
		TOKEN_ERROR,
	})
}

func TestComplexExpression(t *testing.T) {
	source := `
var x = 10;
print (x + 2) * 3 >= 5;
`

	s := New(source)

	tokens := collectTokens(s)

	assertTokenTypes(t, tokens, []TokenType{
		TOKEN_VAR,
		TOKEN_IDENTIFIER,
		TOKEN_EQUAL,
		TOKEN_NUMBER,
		TOKEN_SEMICOLON,

		TOKEN_PRINT,
		TOKEN_LEFT_PAREN,
		TOKEN_IDENTIFIER,
		TOKEN_PLUS,
		TOKEN_NUMBER,
		TOKEN_RIGHT_PAREN,
		TOKEN_STAR,
		TOKEN_NUMBER,
		TOKEN_GREATER_EQUAL,
		TOKEN_NUMBER,
		TOKEN_SEMICOLON,

		TOKEN_EOF,
	})
}
