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

func NewWordRef(pos, length int, dir WordDirection, crossword *Crossword) WordRef {
	return WordRef{
		Pos:       pos,
		Length:    length,
		Dir:       dir,
		Crossword: crossword,
	}
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

type WordRefs []WordRef

func (w WordRefs) Sorted() WordRefs {
	sort.Slice(w, func(i, j int) bool {
		return w[i].Length > w[j].Length
	})
	return w
}
