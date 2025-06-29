package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ahboujelben/go-crossword/cli/renderer"
	"github.com/ahboujelben/go-crossword/modules/cluer"
)

// parseResult holds the parsed command-line arguments
type parseResult struct {
	Rows          int
	Cols          int
	CrosswordSeed int64
	WithClues     bool
	CluesSeed     int64
	Cryptic       bool
	Unsolved      bool
	OllamaModel   string
	OllamaUrl     string
	Threads       int
	Renderer      renderer.Renderer
}

// parseArguments parses command-line arguments and returns a ParseResult
func parseArguments() (*parseResult, error) {
	rows := flag.Int("rows", 13, "number of rows in the crossword ([3, 15])")
	cols := flag.Int("cols", 13, "number of columns in the crossword ([3, 15])")
	crosswordSeed := flag.Int64("crossword-seed", 0, "seed for the crossword generation ([0, 2^63-1], 0 for a random seed)")
	withClues := flag.Bool("with-clues", false, "generate clues using an Ollama model (requires a running Ollama server)")
	cluesSeed := flag.Int64("clues-seed", 0, "seed for the clues generation ([0, 2^63-1], 0 for random seed)")
	cryptic := flag.Bool("cryptic", false, "generate cryptic clues")
	unsolved := flag.Bool("unsolved", false, "hide the crossword solution")
	defaultOllamaUrl := "http://localhost:11434"
	if envUrl := os.Getenv("OLLAMA_URL"); envUrl != "" {
		defaultOllamaUrl = envUrl
	}
	ollamaUrl := flag.String("ollama-url", defaultOllamaUrl, "URL of the Ollama server")
	ollamaModel := flag.String("ollama-model", "llama3:8b", "Ollama model to use")
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

	return &parseResult{
		Rows:          *rows,
		Cols:          *cols,
		Unsolved:      *unsolved,
		WithClues:     *withClues,
		CrosswordSeed: *crosswordSeed,
		CluesSeed:     *cluesSeed,
		Cryptic:       *cryptic,
		OllamaModel:   *ollamaModel,
		OllamaUrl:     *ollamaUrl,
		Threads:       *threads,
		Renderer:      render,
	}, nil
}

// isSeedValid checks if a seed value is valid
func isSeedValid(seed int64) bool {
	return seed >= 0
}

// isSizeValid checks if a crossword size is valid
func isSizeValid(size int) bool {
	return size >= 3 && size <= 15
}
