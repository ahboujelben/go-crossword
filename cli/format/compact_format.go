package format

import (
	"fmt"

	"github.com/ahboujelben/crossword/generator"
)

func CompactFormat(c *generator.Crossword) {
	fmt.Println()
	for letter := range getFormattedLetters(c) {
		fmt.Print(letter)
	}
	fmt.Println()

	fmt.Println("Rows:")
	for word := range getFormattedRowWords(c) {
		fmt.Printf("%-6s%s\n", word[0], word[1])
	}
	fmt.Println()

	fmt.Println("Columns:")
	for word := range getFormattedColumnWords(c) {
		fmt.Printf("%-6s%s\n", word[0], word[1])
	}
	fmt.Println()
}

func getFormattedLetters(c *generator.Crossword) chan string {
	ch := make(chan string)
	go func() {
		for letter := generator.CrosswordLetter(c); letter != nil; letter = letter.Next() {
			switch {
			case letter.IsBlank():
				ch <- "█ "
			case letter.IsEmpty():
				ch <- "  "
			default:
				ch <- string(letter.GetValue()+'A'-'a') + " "
			}
			if letter.Column() == c.Columns()-1 {
				ch <- "\n"
			}
		}
		close(ch)
	}()
	return ch
}
