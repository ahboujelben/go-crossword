package main

import (
	"fmt"

	"github.com/ahboujelben/go-crossword/modules/crossword"
	"github.com/ahboujelben/go-crossword/modules/dictionary"
)

func generateCrossword(parseResult *parseResult) {
	fmt.Println("Generating crossword...")
	crosswordResult := crossword.NewCrossword(crossword.CrosswordConfig{
		Rows:     parseResult.Rows,
		Cols:     parseResult.Cols,
		Seed:     parseResult.CrosswordSeed,
		Threads:  parseResult.Threads,
		WordDict: dictionary.NewWordDictionary(),
	})

	fmt.Printf("\n%s\n\n", parseResult.Renderer.RenderCrossword(crosswordResult.Crossword, true))
	fmt.Println("Crossword generated successfully!")
	fmt.Printf("Seed: %d\n", crosswordResult.Seed)
}
