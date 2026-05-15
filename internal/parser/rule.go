package parser

import (
	"golox/internal/scanner"
)

type ParseFn func()

type ParserRule struct {
	Prefix     ParseFn
	Infix      ParseFn
	Precedence Precedence
}

func (p *Parser) InitParserRules() {
	p.rules = map[scanner.TokenType]ParserRule{
		scanner.TOKEN_LEFT_PAREN:    {p.grouping, nil, PREC_NONE},
		scanner.TOKEN_RIGHT_PAREN:   {nil, nil, PREC_NONE},
		scanner.TOKEN_LEFT_BRACE:    {nil, nil, PREC_NONE},
		scanner.TOKEN_RIGHT_BRACE:   {nil, nil, PREC_NONE},
		scanner.TOKEN_COMMA:         {nil, nil, PREC_NONE},
		scanner.TOKEN_DOT:           {nil, nil, PREC_NONE},
		scanner.TOKEN_MINUS:         {p.unary, p.binary, PREC_TERM},
		scanner.TOKEN_PLUS:          {nil, p.binary, PREC_TERM},
		scanner.TOKEN_SEMICOLON:     {nil, nil, PREC_NONE},
		scanner.TOKEN_SLASH:         {nil, p.binary, PREC_FACTOR},
		scanner.TOKEN_STAR:          {nil, p.binary, PREC_FACTOR},
		scanner.TOKEN_BANG:          {p.unary, nil, PREC_NONE},
		scanner.TOKEN_BANG_EQUAL:    {nil, p.binary, PREC_EQUALITY},
		scanner.TOKEN_EQUAL:         {nil, nil, PREC_NONE},
		scanner.TOKEN_EQUAL_EQUAL:   {nil, p.binary, PREC_COMPARISON},
		scanner.TOKEN_GREATER:       {nil, p.binary, PREC_COMPARISON},
		scanner.TOKEN_GREATER_EQUAL: {nil, p.binary, PREC_COMPARISON},
		scanner.TOKEN_LESS:          {nil, p.binary, PREC_COMPARISON},
		scanner.TOKEN_LESS_EQUAL:    {nil, p.binary, PREC_COMPARISON},
		scanner.TOKEN_IDENTIFIER:    {nil, nil, PREC_NONE},
		scanner.TOKEN_STRING:        {p.string, nil, PREC_NONE},
		scanner.TOKEN_NUMBER:        {p.number, nil, PREC_NONE},
		scanner.TOKEN_AND:           {nil, nil, PREC_NONE},
		scanner.TOKEN_CLASS:         {nil, nil, PREC_NONE},
		scanner.TOKEN_ELSE:          {nil, nil, PREC_NONE},
		scanner.TOKEN_FALSE:         {p.literal, nil, PREC_NONE},
		scanner.TOKEN_FOR:           {nil, nil, PREC_NONE},
		scanner.TOKEN_FUN:           {nil, nil, PREC_NONE},
		scanner.TOKEN_IF:            {nil, nil, PREC_NONE},
		scanner.TOKEN_NIL:           {p.literal, nil, PREC_NONE},
		scanner.TOKEN_OR:            {nil, nil, PREC_NONE},
		scanner.TOKEN_PRINT:         {nil, nil, PREC_NONE},
		scanner.TOKEN_RETURN:        {nil, nil, PREC_NONE},
		scanner.TOKEN_SUPER:         {nil, nil, PREC_NONE},
		scanner.TOKEN_THIS:          {nil, nil, PREC_NONE},
		scanner.TOKEN_TRUE:          {p.literal, nil, PREC_NONE},
		scanner.TOKEN_VAR:           {nil, nil, PREC_NONE},
		scanner.TOKEN_WHILE:         {nil, nil, PREC_NONE},
		scanner.TOKEN_ERROR:         {nil, nil, PREC_NONE},
		scanner.TOKEN_EOF:           {nil, nil, PREC_NONE},
	}
}
