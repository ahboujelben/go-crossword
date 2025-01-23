package generator

import (
	"fmt"
)

type WordRef struct {
	Pos    int
	Length int
	Dir    WordDirection

	Crossword *Crossword
}

func NewHorizontalWordRef(pos int, c *Crossword) *WordRef {
	wordStart := -1
	wordLength := 0
	i := pos
	for i < len(c.Data) {
		if c.Data[i] != '.' {
			if wordStart == -1 {
				wordStart = i
			}
			wordLength++
			if (i+1)%c.Width == 0 || c.Data[i+1] == '.' {
				if wordLength > 1 {
					return &WordRef{
						Pos:       wordStart,
						Length:    wordLength,
						Dir:       Horizontal,
						Crossword: c,
					}
				}
				wordStart = -1
				wordLength = 0
			}
		}
		i++
	}
	return nil
}

func NewVerticalWordRef(pos int, c *Crossword) *WordRef {
	wordStart := -1
	wordLength := 0
	i := pos
	for i < len(c.Data) {
		if c.Data[i] != '.' {
			if wordStart == -1 {
				wordStart = i
			}
			wordLength++
			if i+c.Width >= len(c.Data) || c.Data[i+c.Width] == '.' {
				if wordLength > 1 {
					return &WordRef{
						Pos:       wordStart,
						Length:    wordLength,
						Dir:       Vertical,
						Crossword: c,
					}
				}
				wordStart = -1
				wordLength = 0
			}
		}
		if i == len(c.Data)-1 {
			break
		} else if i+c.Width >= len(c.Data) {
			i = i%c.Width + 1
		} else {
			i += c.Width
		}
	}
	return nil
}

func (w *WordRef) Next() *WordRef {
	if w.Dir == Horizontal {
		nextHorizontalWord := NewHorizontalWordRef(w.Pos+(w.Length-1), w.Crossword)
		if nextHorizontalWord != nil {
			return nextHorizontalWord
		}
		return NewVerticalWordRef(0, w.Crossword)
	}
	return NewVerticalWordRef(w.Pos+(w.Length-1)*w.Crossword.Width, w.Crossword)
}

func (w *WordRef) FirstLetter() *WordLetterRef {
	return NewWordLetterRef(w)
}

func (w *WordRef) IsFilled() bool {
	for letter := w.FirstLetter(); letter != nil; letter = letter.Next() {
		if !letter.IsFilled() {
			return false
		}
	}
	return true
}

func (w *WordRef) GetValue() []byte {
	word := []byte{}
	for letter := w.FirstLetter(); letter != nil; letter = letter.Next() {
		word = append(word, letter.GetValue())
	}
	return word
}

func (w *WordRef) SetValue(value []byte) {
	if len(value) != w.Length {
		panic(fmt.Sprintf("expected length %d but got %d", w.Length, len(value)))
	}
	index := 0
	for letter := w.FirstLetter(); letter != nil; letter = letter.Next() {
		letter.SetValue(value[index])
		index++
	}
}
