package format

import (
	"fmt"

	"github.com/ahboujelben/crossword/generator"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var blackColor = lipgloss.Color("#000")
var whiteColor = lipgloss.Color("#fff")

func StandardFormat(c *generator.Crossword) {
	borderColor := blackColor
	if lipgloss.HasDarkBackground() {
		borderColor = whiteColor
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(borderColor)).
		BorderRow(true).
		Data(newCrosswordCharmWrapper(c))

	fmt.Println(t.Render())
}

type crowssordCharmWrapper struct {
	*generator.Crossword
}

func newCrosswordCharmWrapper(c *generator.Crossword) *crowssordCharmWrapper {
	return &crowssordCharmWrapper{
		Crossword: c,
	}
}

func (w *crowssordCharmWrapper) At(row, column int) string {
	letter := generator.CrosswordLetterAt(w.Crossword, row, column).GetValue()
	switch letter {
	case generator.Blank:
		return "▐█▌"
	case 0:
		return "   "
	default:
		return fmt.Sprintf(" %c ", letter+'A'-'a')
	}
}
