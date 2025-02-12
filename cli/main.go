package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/ahboujelben/crossword/cli/cluer"
	"github.com/ahboujelben/crossword/cli/renderer"
	"github.com/ahboujelben/crossword/generator"
)

func main() {
	parseResult, err := parseArguments()
	if err != nil {
		return
	}

	wordDict := generator.NewWordDict()
	crosswordResult := generator.NewCrossword(generator.CrosswordConfig{
		Rows:     parseResult.rows,
		Cols:     parseResult.cols,
		Seed:     parseResult.crosswordSeed,
		Threads:  parseResult.threads,
		WordDict: wordDict,
	})
	cluesResult := cluer.MakeClues(crosswordResult.Crossword, parseResult.cluesSeed)
	fmt.Println(parseResult.renderer.RenderCrosswordAndClues(crosswordResult.Crossword, cluesResult.Clues))
	fmt.Printf("%-15s: %d\n", "Crossword seed", crosswordResult.Seed)
	fmt.Printf("%-15s: %d\n", "Clues seed", cluesResult.Seed)
}

type ParseResult struct {
	rows          int
	cols          int
	crosswordSeed int64
	cluesSeed     int64
	threads       int
	renderer      renderer.Renderer
}

func parseArguments() (*ParseResult, error) {
	rows := flag.Int("rows", 13, "number of rows in the crossword ([3, 15])")
	cols := flag.Int("cols", 13, "number of columns in the crossword ([3, 15])")
	crosswordSeed := flag.Int64("crossword-seed", 0, "seed for the crossword generation ([0, 2^63-1], 0 for random seed)")
	cluesSeed := flag.Int64("clues-seed", 0, "seed for the clues generation ([0, 2^63-1], 0 for random seed)")
	threads := flag.Int("threads", 100, "number of goroutines to use (>= 1)")
	compact := flag.Bool("compact", false, "compact rendering")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  -rows\n\t%s (default %d)\n", "number of rows in the crossword ([3, 15])", 13)
		fmt.Fprintf(os.Stderr, "  -cols\n\t%s (default %d)\n", "number of columns in the crossword ([3, 15])", 13)
		fmt.Fprintf(os.Stderr, "  -crossword-seed\n\t%s (default %d)\n", "seed for the crossword generation ([0, 2^63-1], 0 for random seed)", 0)
		fmt.Fprintf(os.Stderr, "  -clues-seed\n\t%s (default %d)\n", "seed for the clues generation ([0, 2^63-1], 0 for random seed)", 0)
		fmt.Fprintf(os.Stderr, "  -threads\n\t%s (default %d)\n", "number of goroutines to use (>= 1)", 100)
		fmt.Fprintf(os.Stderr, "  -compact\n\t%s (default %s)\n", "compact rendering", "false")
	}
	flag.Parse()

	if !isDimensionValid(*rows) || !isDimensionValid(*cols) {
		flag.Usage()
		return nil, fmt.Errorf("invalid dimensions")
	}

	if !isSeedValid(*crosswordSeed) || !isSeedValid(*cluesSeed) {
		flag.Usage()
		return nil, fmt.Errorf("invalid seed")
	}

	if *threads < 1 {
		flag.Usage()
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	var render renderer.Renderer = renderer.NewStandardRenderer()
	if *compact {
		render = renderer.NewCompactRenderer()
	}

	return &ParseResult{
		rows:          *rows,
		cols:          *cols,
		crosswordSeed: *crosswordSeed,
		cluesSeed:     *cluesSeed,
		threads:       *threads,
		renderer:      render,
	}, nil
}

func isSeedValid(seed int64) bool {
	return seed >= 0
}

func isDimensionValid(size int) bool {
	return size >= 3 && size <= 15
}
