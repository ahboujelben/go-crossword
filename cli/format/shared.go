package format

import (
	"fmt"

	"github.com/ahboujelben/crossword/generator"
)

func getFormattedRowWords(c *generator.Crossword) chan []string {
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
			ch <- []string{fmt.Sprintf("%d.%d", row+1, currentWordInRow+1), string(word.GetValue())}
			currentWordInRow++
		}
		close(ch)
	}()
	return ch
}
func getFormattedColumnWords(c *generator.Crossword) chan []string {
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
			ch <- []string{fmt.Sprintf("%d.%d", column+1, currentWordInColumn+1), string(word.GetValue())}
			currentWordInColumn++
		}
		close(ch)
	}()
	return ch
}
