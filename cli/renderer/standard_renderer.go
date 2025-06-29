package renderer

import (
	"fmt"
	"os"

	"github.com/ahboujelben/go-crossword/modules/crossword"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
)

var blackColor = lipgloss.Color("#000")
var whiteColor = lipgloss.Color("#fff")

// getBorderTable returns a table with borders
func getBorderTable() *table.Table {
	borderColor := blackColor
	if lipgloss.HasDarkBackground() {
		borderColor = whiteColor
	}

	return table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(borderColor))
}

// StandardRenderer implements a standard rendering of crosswords
type StandardRenderer struct {
}

// NewStandardRenderer creates a new StandardRenderer
func NewStandardRenderer() StandardRenderer {
	return StandardRenderer{}
}

// RenderCrosswordAndClues renders both the crossword and clues
func (f StandardRenderer) RenderCrosswordAndClues(c *crossword.Crossword, clues map[string]string, solved bool) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		f.RenderCrossword(c, solved),
		f.RenderClues(c, clues, solved),
	)
}

// RenderCrossword renders just the crossword grid
func (f StandardRenderer) RenderCrossword(c *crossword.Crossword, solved bool) string {
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

// RenderClues renders just the clues
func (f StandardRenderer) RenderClues(c *crossword.Crossword, clues map[string]string, solved bool) string {
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

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// crosswordCharmWrapper wraps a crossword for use with the charmbracelet table
type crosswordCharmWrapper struct {
	*crossword.Crossword
	solved bool
}

// newCrosswordCharmWrapper creates a new crosswordCharmWrapper
func newCrosswordCharmWrapper(c *crossword.Crossword, solved bool) *crosswordCharmWrapper {
	return &crosswordCharmWrapper{
		Crossword: c,
		solved:    solved,
	}
}

// Columns returns the number of columns in the wrapped crossword
func (w *crosswordCharmWrapper) Columns() int {
	return w.Crossword.Columns() + 1
}

// Rows returns the number of rows in the wrapped crossword
func (w *crosswordCharmWrapper) Rows() int {
	return w.Crossword.Rows() + 1
}

// At returns the content at the given row and column
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
	letter := crossword.CrosswordLetterAt(w.Crossword, row-1, column-1)
	switch {
	case letter.IsBlank():
		return "▐█▌"
	case letter.IsEmpty() || !w.solved:
		return "   "
	default:
		return fmt.Sprintf(" %c ", letter.GetValue()+'A'-'a')
	}
}

// rowsDescriptionWrapper wraps row descriptions for use with the charmbracelet table
type rowsDescriptionWrapper struct {
	words [][]string
}

// newRowsDescriptionWrapper creates a new rowsDescriptionWrapper
func newRowsDescriptionWrapper(c *crossword.Crossword, clues map[string]string, solved bool) *rowsDescriptionWrapper {
	words := [][]string{}
	for word := range getRenderedRowLines(c, clues, solved) {
		words = append(words, word)
	}
	return &rowsDescriptionWrapper{
		words: words,
	}
}

// Columns returns the number of columns in the wrapped description
func (w *rowsDescriptionWrapper) Columns() int {
	return 2
}

// Rows returns the number of rows in the wrapped description
func (w *rowsDescriptionWrapper) Rows() int {
	return len(w.words)
}

// At returns the content at the given row and column
func (w *rowsDescriptionWrapper) At(row, column int) string {
	return w.words[row][column]
}

// columnsDescriptionWrapper wraps column descriptions for use with the charmbracelet table
type columnsDescriptionWrapper struct {
	words [][]string
}

// newColumnsDescriptionWrapper creates a new columnsDescriptionWrapper
func newColumnsDescriptionWrapper(c *crossword.Crossword, clues map[string]string, solved bool) *columnsDescriptionWrapper {
	words := [][]string{}
	for word := range getRenderedColumnLines(c, clues, solved) {
		words = append(words, word)
	}
	return &columnsDescriptionWrapper{
		words: words,
	}
}

// Columns returns the number of columns in the wrapped description
func (w *columnsDescriptionWrapper) Columns() int {
	return 2
}

// Rows returns the number of rows in the wrapped description
func (w *columnsDescriptionWrapper) Rows() int {
	return len(w.words)
}

// At returns the content at the given row and column
func (w *columnsDescriptionWrapper) At(row, column int) string {
	return w.words[row][column]
}
