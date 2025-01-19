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

func NewCrosswordLetterRef(pos int, crossword *Crossword) *CrosswordLetterRef {
	return &CrosswordLetterRef{
		LetterRef: LetterRef{
			Pos:       pos,
			Crossword: crossword,
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
		if l.Pos < l.Word.Pos+l.Word.Length-1 {
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

	if l.Pos/l.Crossword.Width < l.Word.Pos/l.Crossword.Width+l.Word.Length-1 {
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
