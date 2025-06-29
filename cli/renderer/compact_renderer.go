package renderer

import (
	"fmt"

	"github.com/ahboujelben/go-crossword/modules/crossword"
)

// CompactRenderer implements a compact rendering of crosswords
type CompactRenderer struct {
}

// NewCompactRenderer creates a new CompactRenderer
func NewCompactRenderer() CompactRenderer {
	return CompactRenderer{}
}

// RenderCrosswordAndClues renders both the crossword and clues
func (f CompactRenderer) RenderCrosswordAndClues(c *crossword.Crossword, clues map[string]string, solved bool) string {
	result := f.RenderCrossword(c, solved)
	result += f.RenderClues(c, clues, solved)
	return result
}

// RenderCrossword renders just the crossword grid
func (f CompactRenderer) RenderCrossword(c *crossword.Crossword, solved bool) string {
	var result string
	for letter := range getFormattedLetters(c, solved) {
		result += letter
	}

	return result
}

// RenderClues renders just the clues
func (f CompactRenderer) RenderClues(c *crossword.Crossword, clues map[string]string, solved bool) string {
	result := "\n\n"
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

// getFormattedLetters returns a channel of formatted letters for the compact renderer
func getFormattedLetters(c *crossword.Crossword, solved bool) chan string {
	ch := make(chan string)
	go func() {
		for letter := crossword.CrosswordLetter(c); letter != nil; letter = letter.Next() {
			switch {
			case letter.IsBlank():
				ch <- "█ "
			case letter.IsEmpty() || !solved:
				ch <- ". "
			default:
				ch <- string(letter.GetValue()+'A'-'a') + " "
			}
			if letter.Column() == c.Columns()-1 && letter.Row() != c.Rows()-1 {
				ch <- "\n"
			}
		}
		close(ch)
	}()
	return ch
}
