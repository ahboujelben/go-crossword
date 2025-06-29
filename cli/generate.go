package main

import (
	"fmt"

	"github.com/ahboujelben/go-crossword/modules/cluer"
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

	if !parseResult.WithClues {
		fmt.Printf("\n%s\n\n", parseResult.Renderer.RenderCrossword(crosswordResult.Crossword, true))
		fmt.Println("To generate clues for this crossword, run the previous command with these flags:")
		fmt.Printf("  -with-clues -crossword-seed=%d\n", crosswordResult.Seed)
		return
	}

	fmt.Print("Generating clues...\n\n")
	cluesResult := cluer.NewClues(crosswordResult.Crossword, cluer.CluesConfig{
		Seed:        parseResult.CluesSeed,
		Cryptic:     parseResult.Cryptic,
		OllamaModel: parseResult.OllamaModel,
		OllamaUrl:   parseResult.OllamaUrl,
	})

	fmt.Println(parseResult.Renderer.RenderCrosswordAndClues(crosswordResult.Crossword, cluesResult.Clues, !parseResult.Unsolved))
	fmt.Println()
	if parseResult.Unsolved {
		fmt.Println("To reveal the solution for this crossword, run the previous command without the -unsolved flag and with these flags:")
	} else {
		fmt.Println("To display this crossword without the solution, run the previous command with the -unsolved flag and with these flags:")
	}
	fmt.Printf("  -crossword-seed=%d -clues-seed=%d\n\n", crosswordResult.Seed, cluesResult.Seed)
}
