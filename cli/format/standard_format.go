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

	getBorderTable := func() *table.Table {
		return table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(borderColor))
	}

	getDescriptionTable := func(label string) *table.Table {
		return getBorderTable().
			Headers(label).
			StyleFunc(func(row, col int) lipgloss.Style {
				return lipgloss.NewStyle().Padding(0, 1)
			})
	}

	crosswordGrid := getBorderTable().
		BorderRow(true).
		Data(newCrosswordCharmWrapper(c))

	rows := getDescriptionTable("Rows").
		Data(newRowsDescriptionWrapper(c))

	columns := getDescriptionTable("Cols").
		Data(newColumnsDescriptionWrapper(c))

	fmt.Println(lipgloss.JoinHorizontal(
		lipgloss.Top,
		crosswordGrid.Render(),
		lipgloss.NewStyle().MarginLeft(2).Render(rows.Render()),
		lipgloss.NewStyle().MarginLeft(2).Render(columns.Render()),
	))
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
	*generator.Crossword
	words [][]string
}

func newRowsDescriptionWrapper(c *generator.Crossword) *rowsDescriptionWrapper {
	words := [][]string{}
	for word := range getFormattedRowWords(c) {
		words = append(words, word)
	}
	return &rowsDescriptionWrapper{
		Crossword: c,
		words:     words,
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
	*generator.Crossword
	words [][]string
}

func newColumnsDescriptionWrapper(c *generator.Crossword) *columnsDescriptionWrapper {
	words := [][]string{}
	for word := range getFormattedColumnWords(c) {
		words = append(words, word)
	}
	return &columnsDescriptionWrapper{
		Crossword: c,
		words:     words,
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
