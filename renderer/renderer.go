package renderer

import "github.com/ahboujelben/go-crossword/generator"

type Renderer interface {
	RenderCrosswordAndClues(c *generator.Crossword, clues map[string]string, solved bool) string
	RenderCrossword(c *generator.Crossword, solved bool) string
	RenderClues(c *generator.Crossword, clues map[string]string, solved bool) string
}
