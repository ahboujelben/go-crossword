package renderer

import "github.com/ahboujelben/crossword/generator"

type Renderer interface {
	RenderCrosswordAndClues(c *generator.Crossword, clues map[string]string) string
	RenderCrossword(c *generator.Crossword) string
	RenderClues(c *generator.Crossword, clues map[string]string) string
}
