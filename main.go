package main

import (
	"flag"
	"fmt"

	"github.com/ahboujelben/go-crossword/cluer"
	"github.com/ahboujelben/go-crossword/generator"
	"github.com/ahboujelben/go-crossword/renderer"
)

func main() {
	parseResult, err := parseArguments()
	if err != nil {
		fmt.Println(fmt.Errorf("Something is not right: %w", err))
		return
	}

	wordDict := generator.NewWordDict()

	fmt.Println("Generating crossword...")
	crosswordResult := generator.NewCrossword(generator.CrosswordConfig{
		Rows:     parseResult.rows,
		Cols:     parseResult.cols,
		Seed:     parseResult.crosswordSeed,
		Threads:  parseResult.threads,
		WordDict: wordDict,
	})

	if parseResult.omitClues {
		fmt.Printf("%s\n", parseResult.renderer.RenderCrossword(crosswordResult.Crossword, true))
		fmt.Printf("%-15s: %d\n", "Crossword seed", crosswordResult.Seed)
		return
	}

	fmt.Println("Generating clues...")
	cluesResult := cluer.NewClues(crosswordResult.Crossword, cluer.CluesConfig{
		Seed:        parseResult.cluesSeed,
		Difficulty:  parseResult.cluesDifficulty,
		OllamaModel: parseResult.ollamaModel,
		OllamaUrl:   parseResult.ollamaUrl,
	})

	fmt.Println()
	fmt.Printf("%-15s: %d\n", "Crossword seed", crosswordResult.Seed)
	fmt.Printf("%-15s: %d\n", "Clues seed", cluesResult.Seed)
	println()

	fmt.Println(parseResult.renderer.RenderCrosswordAndClues(crosswordResult.Crossword, cluesResult.Clues, parseResult.solved))
}

type ParseResult struct {
	rows            int
	cols            int
	crosswordSeed   int64
	omitClues       bool
	cluesDifficulty cluer.ClueDifficulty
	solved          bool
	cluesSeed       int64
	ollamaModel     string
	ollamaUrl       string
	threads         int
	renderer        renderer.Renderer
}

func parseArguments() (*ParseResult, error) {
	rows := flag.Int("rows", 13, "number of rows in the crossword ([3, 15])")
	cols := flag.Int("cols", 13, "number of columns in the crossword ([3, 15])")
	crosswordSeed := flag.Int64("crossword-seed", 0, "seed for the crossword generation ([0, 2^63-1], 0 for a random seed)")
	omitClues := flag.Bool("omit-clues", false, "generate only the crossword without the clues")
	cryptic := flag.Bool("cryptic", false, "generate cryptic clues")
	solved := flag.Bool("solved", false, "show the solved crossword")
	cluesSeed := flag.Int64("clues-seed", 0, "seed for the clues generation ([0, 2^63-1], 0 for random seed)")
	ollamaUrl := flag.String("ollama-url", "http://localhost:11434", "URL of the Ollama server")
	ollamaModel := flag.String("ollama-model", "llama3.1:8b", "Ollama model to use")
	threads := flag.Int("threads", 100, "number of goroutines to use (>= 1)")
	compact := flag.Bool("compact", false, "compact rendering")

	flag.Parse()

	if !isDimensionValid(*rows) || !isDimensionValid(*cols) {
		return nil, fmt.Errorf("invalid dimensions")
	}

	if !*omitClues {
		if err := cluer.CheckOllama(*ollamaUrl, *ollamaModel); err != nil {
			return nil, err
		}
	}

	if !isSeedValid(*crosswordSeed) || !isSeedValid(*cluesSeed) {
		return nil, fmt.Errorf("invalid seed")
	}

	if *threads < 1 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	cluesDifficulty := cluer.ClueDifficultyNormal
	if *cryptic {
		cluesDifficulty = cluer.ClueDifficultyCryptic
	}

	var render renderer.Renderer = renderer.NewStandardRenderer()
	if *compact {
		render = renderer.NewCompactRenderer()
	}

	return &ParseResult{
		rows:            *rows,
		cols:            *cols,
		solved:          *solved,
		omitClues:       *omitClues,
		crosswordSeed:   *crosswordSeed,
		cluesSeed:       *cluesSeed,
		cluesDifficulty: cluesDifficulty,
		ollamaModel:     *ollamaModel,
		ollamaUrl:       *ollamaUrl,
		threads:         *threads,
		renderer:        render,
	}, nil
}

func isSeedValid(seed int64) bool {
	return seed >= 0
}

func isDimensionValid(size int) bool {
	return size >= 3 && size <= 15
}
