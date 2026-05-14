package scanner

import "fmt"

type Scanner struct {
	Start   int
	Current int
	Line    int

	source       string
	sourceLength int
}

func New(source string) *Scanner {
	src := source
	return &Scanner{
		Start:        0,
		Current:      0,
		Line:         1,
		source:       src,
		sourceLength: len(src),
	}
}

func (s *Scanner) ScanTokens() {
	line := -1
	for {
		token := s.ScanToken()
		if token.Line != line {
			fmt.Printf("%4d ", token.Line)
			line = token.Line
		} else {
			fmt.Printf("   | ")
		}
		fmt.Printf("type: %2d, length: %d, lexeme: '%s'\n", token.Type, token.Length, token.Lexeme)

		if token.Type == TOKEN_EOF {
			break
		}
	}
}

func (s *Scanner) Free() {
	s.source = ""
	s.sourceLength = 0
	s.Start = 0
	s.Current = 0
	s.Line = 0
}
