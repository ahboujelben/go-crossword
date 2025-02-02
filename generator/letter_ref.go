package generator

type LetterRef struct {
	Pos       int
	crossword *Crossword
}

func (l *LetterRef) GetValue() byte {
	return l.crossword.data[l.Pos]
}

func (l *LetterRef) SetValue(value byte) {
	l.crossword.data[l.Pos] = value
}

func (l *LetterRef) IsEmpty() bool {
	return l.crossword.data[l.Pos] == 0
}

type CrosswordLetterRef struct {
	LetterRef
}

func CrosswordLetter(c *Crossword) *CrosswordLetterRef {
	if len(c.data) == 0 {
		return nil
	}

	return &CrosswordLetterRef{
		LetterRef: LetterRef{
			crossword: c,
		},
	}
}

func CrosswordLetterAt(c *Crossword, row, column int) *CrosswordLetterRef {
	return &CrosswordLetterRef{
		LetterRef: LetterRef{
			Pos:       row*c.columns + column,
			crossword: c,
		},
	}
}

func (l *CrosswordLetterRef) Next() *CrosswordLetterRef {
	if l.Pos+1 < len(l.crossword.data) {
		return &CrosswordLetterRef{
			LetterRef: LetterRef{
				Pos:       l.Pos + 1,
				crossword: l.crossword,
			},
		}
	}
	return nil
}

type WordLetterRef struct {
	LetterRef
	word *WordRef
}

func WordLetter(word *WordRef) *WordLetterRef {
	return &WordLetterRef{
		LetterRef: LetterRef{
			Pos:       word.Pos,
			crossword: word.crossword,
		},
		word: word,
	}
}

func (l *WordLetterRef) Next() *WordLetterRef {
	if l.word.Dir == Horizontal {
		if l.Pos+1 < l.word.Pos+l.word.Length {
			return &WordLetterRef{
				LetterRef: LetterRef{
					Pos:       l.Pos + 1,
					crossword: l.crossword,
				},
				word: l.word,
			}
		}
		return nil
	}

	if l.Pos+l.crossword.columns < l.word.Pos+l.word.Length*l.crossword.columns {
		return &WordLetterRef{
			LetterRef: LetterRef{
				Pos:       l.Pos + l.crossword.columns,
				crossword: l.crossword,
			},
			word: l.word,
		}
	}
	return nil
}
