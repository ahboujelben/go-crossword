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

func (f CompactRenderer) RenderCrosswordAndClues(c *generator.Crossword, clues map[string]string, solved bool) string {
	result := f.RenderCrossword(c, solved)
	result += f.RenderClues(c, clues, solved)
	return result
}

func (f CompactRenderer) RenderCrossword(c *generator.Crossword, solved bool) string {
	var result string
	for letter := range getFormattedLetters(c, solved) {
		result += letter
	}
	result += "\n"

	return result
}

func (f CompactRenderer) RenderClues(c *generator.Crossword, clues map[string]string, solved bool) string {
	var result string
	result += "Rows:\n"
	for word := range getRenderedRowLines(c, clues, solved) {
		if solved {
			result += fmt.Sprintf("%s\n", word[1])
		} else {
			result += fmt.Sprintf("%-6s%s\n", word[0], word[1])
		}
	}
	result += "\n"

	result += "Columns:\n"
	for word := range getRenderedColumnLines(c, clues, solved) {
		if solved {
			result += fmt.Sprintf("%s\n", word[1])
		} else {
			result += fmt.Sprintf("%-6s%s\n", word[0], word[1])
		}
	}

	return result
}

func getFormattedLetters(c *generator.Crossword, solved bool) chan string {
	ch := make(chan string)
	go func() {
		for letter := generator.CrosswordLetter(c); letter != nil; letter = letter.Next() {
			switch {
			case letter.IsBlank():
				ch <- "█ "
			case letter.IsEmpty() || !solved:
				ch <- ". "
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
