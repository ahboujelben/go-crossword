package renderer

import (
	"github.com/ahboujelben/go-crossword/modules/crossword"
)

// Renderer defines the interface for rendering crosswords
type Renderer interface {
	RenderCrosswordAndClues(c *crossword.Crossword, clues map[string]string, solved bool) string
	RenderCrossword(c *crossword.Crossword, solved bool) string
	RenderClues(c *crossword.Crossword, clues map[string]string, solved bool) string
}
