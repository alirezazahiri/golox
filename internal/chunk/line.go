package chunk

type LineStart struct {
	Line  int
	Count int
}

func (c *Chunk) GetLine(instruction int) int {
	current := 0

	for _, run := range c.Lines {
		current += run.Count

		if instruction < current {
			return run.Line
		}
	}

	return -1
}
