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
		scanner.TOKEN_LEFT_PAREN:    ParserRule{p.grouping, nil, PREC_NONE},
		scanner.TOKEN_RIGHT_PAREN:   ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_LEFT_BRACE:    ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_RIGHT_BRACE:   ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_COMMA:         ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_DOT:           ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_MINUS:         ParserRule{p.unary, p.binary, PREC_TERM},
		scanner.TOKEN_PLUS:          ParserRule{nil, p.binary, PREC_TERM},
		scanner.TOKEN_SEMICOLON:     ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_SLASH:         ParserRule{nil, p.binary, PREC_FACTOR},
		scanner.TOKEN_STAR:          ParserRule{nil, p.binary, PREC_FACTOR},
		scanner.TOKEN_BANG:          ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_BANG_EQUAL:    ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_EQUAL:         ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_EQUAL_EQUAL:   ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_GREATER:       ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_GREATER_EQUAL: ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_LESS:          ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_LESS_EQUAL:    ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_IDENTIFIER:    ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_STRING:        ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_NUMBER:        ParserRule{p.number, nil, PREC_NONE},
		scanner.TOKEN_AND:           ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_CLASS:         ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_ELSE:          ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_FALSE:         ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_FOR:           ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_FUN:           ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_IF:            ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_NIL:           ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_OR:            ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_PRINT:         ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_RETURN:        ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_SUPER:         ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_THIS:          ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_TRUE:          ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_VAR:           ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_WHILE:         ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_ERROR:         ParserRule{nil, nil, PREC_NONE},
		scanner.TOKEN_EOF:           ParserRule{nil, nil, PREC_NONE},
	}
}
