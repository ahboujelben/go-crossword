package renderer

import (
	"fmt"
	"os"

	"github.com/ahboujelben/crossword/generator"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
)

var blackColor = lipgloss.Color("#000")
var whiteColor = lipgloss.Color("#fff")

func getBorderTable() *table.Table {
	borderColor := blackColor
	if lipgloss.HasDarkBackground() {
		borderColor = whiteColor
	}

	return table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(borderColor))
}

type StandardRenderer struct {
}

func NewStandardRenderer() StandardRenderer {
	return StandardRenderer{}
}

func (f StandardRenderer) RenderCrosswordAndClues(c *generator.Crossword, clues map[string]string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		f.RenderCrossword(c),
		f.RenderClues(c, clues),
	)
}

func (f StandardRenderer) RenderCrossword(c *generator.Crossword) string {
	crosswordGrid := getBorderTable().
		BorderRow(true).
		Data(newCrosswordCharmWrapper(c))

	return crosswordGrid.Render()
}

func (f StandardRenderer) RenderClues(c *generator.Crossword, clues map[string]string) string {
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth = 120
	}
	getDescriptionTable := func(label string) *table.Table {
		return getBorderTable().
			Headers(label).
			StyleFunc(func(row, col int) lipgloss.Style {
				s := lipgloss.NewStyle().Padding(0, 1)
				if col == 1 {
					s = s.Width(min((termWidth - 4*c.Columns() - 14), 80))
				}
				return s
			})
	}

	rows := getDescriptionTable("Rows").
		Data(newRowsDescriptionWrapper(c, clues))

	columns := getDescriptionTable("Cols").
		Data(newColumnsDescriptionWrapper(c, clues))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().MarginLeft(2).Render(rows.Render()),
		lipgloss.NewStyle().MarginLeft(2).Render(columns.Render()),
	)
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
	letter := generator.CrosswordLetterAt(w.Crossword, row, column)
	switch {
	case letter.IsBlank():
		return "▐█▌"
	case letter.IsEmpty():
		return "   "
	default:
		return fmt.Sprintf(" %c ", letter.GetValue()+'A'-'a')
	}
}

type rowsDescriptionWrapper struct {
	words [][]string
}

func newRowsDescriptionWrapper(c *generator.Crossword, clues map[string]string) *rowsDescriptionWrapper {
	words := [][]string{}
	for word := range getRenderedRowLines(c, clues) {
		words = append(words, word)
	}
	return &rowsDescriptionWrapper{
		words: words,
	}
}

func (w *rowsDescriptionWrapper) Columns() int {
	return 2
}

func (w *rowsDescriptionWrapper) Rows() int {
	return len(w.words)
}

func (w *rowsDescriptionWrapper) At(row, column int) string {
	return w.words[row][column]
}

type columnsDescriptionWrapper struct {
	words [][]string
}

func newColumnsDescriptionWrapper(c *generator.Crossword, clues map[string]string) *columnsDescriptionWrapper {
	words := [][]string{}
	for word := range getRenderedColumnLines(c, clues) {
		words = append(words, word)
	}
	return &columnsDescriptionWrapper{
		words: words,
	}
}

func (w *columnsDescriptionWrapper) Columns() int {
	return 2
}

func (w *columnsDescriptionWrapper) Rows() int {
	return len(w.words)
}

func (w *columnsDescriptionWrapper) At(row, column int) string {
	return w.words[row][column]
}
