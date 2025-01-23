package generator

type LetterRef struct {
	Pos int

	Crossword *Crossword
}

func (l *LetterRef) IsFilled() bool {
	return l.Crossword.Data[l.Pos] != 0
}

func (l *LetterRef) GetValue() byte {
	return l.Crossword.Data[l.Pos]
}

func (l *LetterRef) SetValue(value byte) {
	l.Crossword.Data[l.Pos] = value
}

type CrosswordLetterRef struct {
	LetterRef
}

func NewCrosswordLetterRef(pos int, c *Crossword) *CrosswordLetterRef {
	if len(c.Data) == 0 {
		return nil
	}

	return &CrosswordLetterRef{
		LetterRef: LetterRef{
			Pos:       pos,
			Crossword: c,
		},
	}
}

func (l *CrosswordLetterRef) Next() *CrosswordLetterRef {
	if l.Pos+1 < len(l.Crossword.Data) {
		return &CrosswordLetterRef{
			LetterRef: LetterRef{
				Pos:       l.Pos + 1,
				Crossword: l.Crossword,
			},
		}
	}
	return nil
}

type WordLetterRef struct {
	LetterRef
	Word *WordRef
}

func NewWordLetterRef(word *WordRef) *WordLetterRef {
	return &WordLetterRef{
		LetterRef: LetterRef{
			Pos:       word.Pos,
			Crossword: word.Crossword,
		},
		Word: word,
	}
}

func (l *WordLetterRef) Next() *WordLetterRef {
	if l.Word.Dir == Horizontal {
		if l.Pos+1 < l.Word.Pos+l.Word.Length {
			return &WordLetterRef{
				LetterRef: LetterRef{
					Pos:       l.Pos + 1,
					Crossword: l.Crossword,
				},
				Word: l.Word,
			}
		}
		return nil
	}

	if l.Pos+l.Crossword.Width < l.Word.Pos+l.Word.Length*l.Crossword.Width {
		return &WordLetterRef{
			LetterRef: LetterRef{
				Pos:       l.Pos + l.Crossword.Width,
				Crossword: l.Crossword,
			},
			Word: l.Word,
		}
	}
	return nil
}
