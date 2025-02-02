package generator

import (
	"fmt"
)

type WordRef struct {
	Pos    int
	Length int
	Dir    WordDirection

	crossword *Crossword
}

type WordDirection int

const (
	Horizontal WordDirection = iota
	Vertical
)

func Word(c *Crossword) *WordRef {
	return horizontalWord(0, c)
}

func horizontalWord(pos int, c *Crossword) *WordRef {
	wordStart := -1
	wordLength := 0
	i := pos
	for i < len(c.data) {
		if c.data[i] != '.' {
			if wordStart == -1 {
				wordStart = i
			}
			wordLength++
			if (i+1)%c.Width == 0 || c.data[i+1] == '.' {
				if wordLength > 1 {
					return &WordRef{
						Pos:       wordStart,
						Length:    wordLength,
						Dir:       Horizontal,
						crossword: c,
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

func verticalWord(pos int, c *Crossword) *WordRef {
	wordStart := -1
	wordLength := 0
	i := pos
	for i < len(c.data) {
		if c.data[i] != '.' {
			if wordStart == -1 {
				wordStart = i
			}
			wordLength++
			if i+c.Width >= len(c.data) || c.data[i+c.Width] == '.' {
				if wordLength > 1 {
					return &WordRef{
						Pos:       wordStart,
						Length:    wordLength,
						Dir:       Vertical,
						crossword: c,
					}
				}
				wordStart = -1
				wordLength = 0
			}
		}
		if i == len(c.data)-1 {
			break
		} else if i+c.Width >= len(c.data) {
			i = i%c.Width + 1
		} else {
			i += c.Width
		}
	}
	return nil
}

func (w *WordRef) Next() *WordRef {
	if w.Dir == Horizontal {
		nextHorizontalWord := horizontalWord(w.Pos+(w.Length-1), w.crossword)
		if nextHorizontalWord != nil {
			return nextHorizontalWord
		}
		return verticalWord(0, w.crossword)
	}
	return verticalWord(w.Pos+(w.Length-1)*w.crossword.Width, w.crossword)
}

func (w *WordRef) GetValue() []byte {
	word := []byte{}
	for letter := WordLetter(w); letter != nil; letter = letter.Next() {
		word = append(word, letter.GetValue())
	}
	return word
}

func (w *WordRef) SetValue(value []byte) {
	if len(value) != w.Length {
		panic(fmt.Sprintf("expected length %d but got %d", w.Length, len(value)))
	}
	index := 0
	for letter := WordLetter(w); letter != nil; letter = letter.Next() {
		letter.SetValue(value[index])
		index++
	}
}

func (w *WordRef) IsFilled() bool {
	for letter := WordLetter(w); letter != nil; letter = letter.Next() {
		if letter.IsEmpty() {
			return false
		}
	}
	return true
}
