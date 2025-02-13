package renderer

import (
	"fmt"
	"os"

	"github.com/ahboujelben/go-crossword/generator"

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

func (f StandardRenderer) RenderCrosswordAndClues(c *generator.Crossword, clues map[string]string, solved bool) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		f.RenderCrossword(c, solved),
		f.RenderClues(c, clues, solved),
	)
}

func (f StandardRenderer) RenderCrossword(c *generator.Crossword, solved bool) string {
	crosswordGrid := getBorderTable().
		BorderRow(true).
		StyleFunc(func(row, col int) lipgloss.Style {
			s := lipgloss.NewStyle()
			if row == 0 || col == 0 {
				s = s.Foreground(lipgloss.Color("#00ff00"))
			}
			return s
		}).
		Data(newCrosswordCharmWrapper(c, solved))

	return crosswordGrid.Render()
}

func (f StandardRenderer) RenderClues(c *generator.Crossword, clues map[string]string, solved bool) string {
	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		termWidth = 120
	}
	getDescriptionTable := func(label string) *table.Table {
		return getBorderTable().
			Headers(label).
			StyleFunc(func(row, col int) lipgloss.Style {
				s := lipgloss.NewStyle().Padding(0, 1)
				if row >= 0 && col == 0 {
					s = s.Foreground(lipgloss.Color("#00ff00"))
				}
				if col == 1 {
					s = s.Width(min((termWidth - 4*c.Columns() - 17), 80))
				}
				return s
			})
	}

	rows := getDescriptionTable("Rows").
		Data(newRowsDescriptionWrapper(c, clues, solved))

	columns := getDescriptionTable("Cols").
		Data(newColumnsDescriptionWrapper(c, clues, solved))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().MarginLeft(2).Render(rows.Render()),
		lipgloss.NewStyle().MarginLeft(2).Render(columns.Render()),
	)
}

type crosswordCharmWrapper struct {
	*generator.Crossword
	solved bool
}

func newCrosswordCharmWrapper(c *generator.Crossword, solved bool) *crosswordCharmWrapper {
	return &crosswordCharmWrapper{
		Crossword: c,
		solved:    solved,
	}
}

func (w *crosswordCharmWrapper) Columns() int {
	return w.Crossword.Columns() + 1
}

func (w *crosswordCharmWrapper) Rows() int {
	return w.Crossword.Rows() + 1
}

func (w *crosswordCharmWrapper) At(row, column int) string {
	if row == 0 && column == 0 {
		return "   "
	}
	if row == 0 {
		return fmt.Sprintf(" %-2d", column)
	}
	if column == 0 {
		return fmt.Sprintf(" %-2d", row)
	}
	letter := generator.CrosswordLetterAt(w.Crossword, row-1, column-1)
	switch {
	case letter.IsBlank():
		return "▐█▌"
	case letter.IsEmpty() || !w.solved:
		return "   "
	default:
		return fmt.Sprintf(" %c ", letter.GetValue()+'A'-'a')
	}
}

type rowsDescriptionWrapper struct {
	words [][]string
}

func newRowsDescriptionWrapper(c *generator.Crossword, clues map[string]string, solved bool) *rowsDescriptionWrapper {
	words := [][]string{}
	for word := range getRenderedRowLines(c, clues, solved) {
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

func newColumnsDescriptionWrapper(c *generator.Crossword, clues map[string]string, solved bool) *columnsDescriptionWrapper {
	words := [][]string{}
	for word := range getRenderedColumnLines(c, clues, solved) {
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
