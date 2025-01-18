package generator

import (
	"fmt"
	"sort"
)

type WordRef struct {
	Pos    int
	Length int
	Dir    WordDirection

	Crossword *Crossword
}

type Word []byte

func NewWordRef(pos, length int, dir WordDirection, crossword *Crossword) WordRef {
	return WordRef{
		Pos:       pos,
		Length:    length,
		Dir:       dir,
		Crossword: crossword,
	}
}

func (w WordRef) Letters() []LetterRef {
	refs := make([]LetterRef, w.Length)
	if w.Dir == Horizontal {
		for i := 0; i < w.Length; i++ {
			refs[i] = LetterRef{
				Pos:       w.Pos + i,
				Crossword: w.Crossword,
			}
		}
	} else {
		for i := 0; i < w.Length; i++ {
			refs[i] = LetterRef{
				Pos:       w.Pos + i*w.Crossword.Width,
				Crossword: w.Crossword,
			}
		}
	}
	return refs
}

func (w WordRef) IsFilled() bool {
	for _, letter := range w.Letters() {
		if !letter.IsFilled() {
			return false
		}
	}
	return true
}

func (w WordRef) GetValue() Word {
	word := make(Word, w.Length)
	for i, letter := range w.Letters() {
		word[i] = letter.GetValue()
	}
	return word
}

func (w WordRef) SetValue(value Word) {
	if len(value) != w.Length {
		panic(fmt.Sprintf("expected length %d but got %d", w.Length, len(value)))
	}
	for i, letter := range w.Letters() {
		letter.SetValue(value[i])
	}
}

type WordRefs []WordRef

func (w WordRefs) Sorted() WordRefs {
	sort.Slice(w, func(i, j int) bool {
		return w[i].Length > w[j].Length
	})
	return w
}
