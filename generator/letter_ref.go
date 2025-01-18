package generator

type LetterRef struct {
	Pos int

	Crossword *Crossword
}

func (l LetterRef) IsFilled() bool {
	return l.Crossword.Data[l.Pos] != 0
}

func (l LetterRef) GetValue() byte {
	return l.Crossword.Data[l.Pos]
}

func (l LetterRef) SetValue(value byte) {
	l.Crossword.Data[l.Pos] = value
}
