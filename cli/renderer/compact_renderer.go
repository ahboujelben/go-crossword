package renderer

import (
	"fmt"

	"github.com/ahboujelben/crossword/generator"
)

type CompactRenderer struct {
}

func NewCompactRenderer() CompactRenderer {
	return CompactRenderer{}
}

func (f CompactRenderer) RenderCrosswordAndClues(c *generator.Crossword, clues map[string]string) string {
	result := f.RenderCrossword(c)
	result += f.RenderClues(c, clues)
	return result
}

func (f CompactRenderer) RenderCrossword(c *generator.Crossword) string {
	var result string
	result += "\n"
	for letter := range getFormattedLetters(c) {
		result += letter
	}
	result += "\n"

	return result
}

func (f CompactRenderer) RenderClues(c *generator.Crossword, clues map[string]string) string {
	var result string
	result += "Rows:\n"
	for word := range getRenderedRowLines(c, clues) {
		result += fmt.Sprintf("%-6s%s\n", word[0], word[1])
	}
	result += "\n"

	result += "Columns:\n"
	for word := range getRenderedColumnLines(c, clues) {
		result += fmt.Sprintf("%-6s%s\n", word[0], word[1])
	}
	result += "\n"

	return result
}

func getFormattedLetters(c *generator.Crossword) chan string {
	ch := make(chan string)
	go func() {
		for letter := generator.CrosswordLetter(c); letter != nil; letter = letter.Next() {
			switch {
			case letter.IsBlank():
				ch <- "█ "
			case letter.IsEmpty():
				ch <- "  "
			default:
				ch <- string(letter.GetValue()+'A'-'a') + " "
			}
			if letter.Column() == c.Columns()-1 {
				ch <- "\n"
			}
		}
		close(ch)
	}()
	return ch
}
