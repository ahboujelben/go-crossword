package renderer

import (
	"fmt"
	"strings"

	"github.com/ahboujelben/crossword/generator"
)

func getRenderedRowLines(c *generator.Crossword, clues map[string]string) chan []string {
	ch := make(chan []string)
	go func() {
		currentRow := 0
		currentWordInRow := 0
		for word := generator.RowWord(c); word != nil; word = word.Next() {
			row := word.Row()
			if row != currentRow {
				currentRow = row
				currentWordInRow = 0
			}
			rowLine := make([]string, 0, 2)
			rowWord := string(word.GetValue())
			rowLine = append(rowLine, rowWord)
			if clue, ok := clues[rowWord]; ok {
				rowLine = append(rowLine, clue)
			}
			rowLineString := strings.Join(rowLine, ": ")
			ch <- []string{fmt.Sprintf("%d.%d", row+1, currentWordInRow+1), rowLineString}
			currentWordInRow++
		}
		close(ch)
	}()
	return ch
}
func getRenderedColumnLines(c *generator.Crossword, clues map[string]string) chan []string {
	ch := make(chan []string)
	go func() {
		currentColumn := 0
		currentWordInColumn := 0
		for word := generator.ColumnWord(c); word != nil; word = word.Next() {
			column := word.Column()
			if column != currentColumn {
				currentColumn = column
				currentWordInColumn = 0
			}
			columnLine := make([]string, 0, 2)
			columnWord := string(word.GetValue())
			columnLine = append(columnLine, columnWord)
			if clue, ok := clues[columnWord]; ok {
				columnLine = append(columnLine, clue)
			}
			columnLineString := strings.Join(columnLine, ": ")
			ch <- []string{fmt.Sprintf("%d.%d", column+1, currentWordInColumn+1), columnLineString}
			currentWordInColumn++
		}
		close(ch)
	}()
	return ch
}
