package main

import (
	"flag"
	"fmt"

	"github.com/ahboujelben/go-crossword/cluer"
	"github.com/ahboujelben/go-crossword/dictionary"
	"github.com/ahboujelben/go-crossword/generator"
	"github.com/ahboujelben/go-crossword/renderer"
)

func main() {
	parseResult, err := parseArguments()
	if err != nil {
		fmt.Println(fmt.Errorf("Something is not right: %w", err))
		return
	}

	fmt.Println("Generating crossword...")
	crosswordResult := generator.NewCrossword(generator.CrosswordConfig{
		Rows:     parseResult.rows,
		Cols:     parseResult.cols,
		Seed:     parseResult.crosswordSeed,
		Threads:  parseResult.threads,
		WordDict: dictionary.NewWordDictionary(),
	})

	if !parseResult.withClues {
		fmt.Printf("%s\n\n", parseResult.renderer.RenderCrossword(crosswordResult.Crossword, true))
		fmt.Println("To generate clues for this crossword, run the previous command with these additional flags:")
		fmt.Printf("  -with-clues -crossword-seed=%d\n\n", crosswordResult.Seed)
		fmt.Printf("Run the program with -h for more information on how to configure Ollama url and model.\n")
		return
	}

	fmt.Print("Generating clues...\n\n")
	cluesResult := cluer.NewClues(crosswordResult.Crossword, cluer.CluesConfig{
		Seed:        parseResult.cluesSeed,
		Cryptic:     parseResult.cryptic,
		OllamaModel: parseResult.ollamaModel,
		OllamaUrl:   parseResult.ollamaUrl,
	})

	fmt.Println(parseResult.renderer.RenderCrosswordAndClues(crosswordResult.Crossword, cluesResult.Clues, !parseResult.unsolved))
	fmt.Println()
	if parseResult.unsolved {
		fmt.Println("To reveal the solution for this crossword, run the previous command without the -unsolved flag and with these flags:")
	} else {
		fmt.Println("To display this crossword without the solution, run the previous command with the -unsolved flag and with these flags:")
	}
	fmt.Printf("  -crossword-seed=%d -clues-seed=%d\n\n", crosswordResult.Seed, cluesResult.Seed)
}

type ParseResult struct {
	rows          int
	cols          int
	crosswordSeed int64
	withClues     bool
	cluesSeed     int64
	cryptic       bool
	unsolved      bool
	ollamaModel   string
	ollamaUrl     string
	threads       int
	renderer      renderer.Renderer
}

func parseArguments() (*ParseResult, error) {
	rows := flag.Int("rows", 13, "number of rows in the crossword ([3, 15])")
	cols := flag.Int("cols", 13, "number of columns in the crossword ([3, 15])")
	crosswordSeed := flag.Int64("crossword-seed", 0, "seed for the crossword generation ([0, 2^63-1], 0 for a random seed)")
	withClues := flag.Bool("with-clues", false, "generate clues using an Ollama model (requires a running Ollama server)")
	cluesSeed := flag.Int64("clues-seed", 0, "seed for the clues generation ([0, 2^63-1], 0 for random seed)")
	cryptic := flag.Bool("cryptic", false, "generate cryptic clues")
	unsolved := flag.Bool("unsolved", false, "hide the crossword solution")
	ollamaUrl := flag.String("ollama-url", "http://localhost:11434", "URL of the Ollama server")
	ollamaModel := flag.String("ollama-model", "llama3.1:8b", "Ollama model to use")
	threads := flag.Int("threads", 100, "number of goroutines to use (>= 1)")
	compact := flag.Bool("compact", false, "compact rendering")

	flag.Parse()

	if !isSizeValid(*rows) || !isSizeValid(*cols) {
		return nil, fmt.Errorf("invalid dimensions")
	}

	if !isSeedValid(*crosswordSeed) {
		return nil, fmt.Errorf("invalid crossword seed")
	}

	if *withClues {
		if err := cluer.CheckOllamaServer(*ollamaUrl, *ollamaModel); err != nil {
			return nil, err
		}
	}

	if !isSeedValid(*cluesSeed) {
		return nil, fmt.Errorf("invalid clues seed")
	}
	if *threads < 1 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	var render renderer.Renderer = renderer.NewStandardRenderer()
	if *compact {
		render = renderer.NewCompactRenderer()
	}

	return &ParseResult{
		rows:          *rows,
		cols:          *cols,
		unsolved:      *unsolved,
		withClues:     *withClues,
		crosswordSeed: *crosswordSeed,
		cluesSeed:     *cluesSeed,
		cryptic:       *cryptic,
		ollamaModel:   *ollamaModel,
		ollamaUrl:     *ollamaUrl,
		threads:       *threads,
		renderer:      render,
	}, nil
}

func isSeedValid(seed int64) bool {
	return seed >= 0
}

func isSizeValid(size int) bool {
	return size >= 3 && size <= 15
}
