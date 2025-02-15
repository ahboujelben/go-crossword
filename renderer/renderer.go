package renderer

import "github.com/ahboujelben/go-crossword/crossword"

type Renderer interface {
	RenderCrosswordAndClues(c *crossword.Crossword, clues map[string]string, solved bool) string
	RenderCrossword(c *crossword.Crossword, solved bool) string
	RenderClues(c *crossword.Crossword, clues map[string]string, solved bool) string
}
