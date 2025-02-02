package format

import (
	"fmt"

	"github.com/ahboujelben/crossword/generator"
)

func CompactFormat(c *generator.Crossword) {
	fmt.Println()
	for letter := generator.CrosswordLetter(c); letter != nil; letter = letter.Next() {
		switch letter.GetValue() {
		case generator.Blank:
			fmt.Print("█ ")
		case 0:
			fmt.Print("  ")
		default:
			fmt.Printf("%c ", letter.GetValue()+'A'-'a')
		}
		if letter.Pos%c.Columns() == c.Columns()-1 {
			fmt.Println()
		}
	}
	fmt.Println()
}
