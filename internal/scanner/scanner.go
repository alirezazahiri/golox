package scanner

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
