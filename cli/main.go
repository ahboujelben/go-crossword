package main

import (
	_ "embed"
	"flag"
	"fmt"

	"github.com/ahboujelben/crossword/cli/format"
	"github.com/ahboujelben/crossword/generator"
)

func main() {
	parseResult, err := parseArguments()
	if err != nil {
		return
	}

	wordDict := generator.NewWordDict()
	crossword := generator.NewCrossword(generator.CrosswordConfig{
		Rows:        parseResult.rows,
		Columns:     parseResult.columns,
		Concurrency: parseResult.concurrency,
		WordDict:    wordDict,
	})
	crossword.Print(parseResult.formatter)
}

type ParseResult struct {
	rows        int
	columns     int
	concurrency int
	formatter   func(c *generator.Crossword)
}

func parseArguments() (*ParseResult, error) {
	rows := flag.Int("rows", 13, "number of rows in the crossword (valid values: [3, 13])")
	columns := flag.Int("cols", 13, "number of columns in the crossword (valid values: [3, 13])")
	concurrency := flag.Int("conc", 100, "number of goroutines to use (valid values: >= 1)")
	isCompact := flag.Bool("compact", false, "prints each letter using one character")
	flag.Parse()

	if !isDimensionValid(*rows) || !isDimensionValid(*columns) {
		flag.Usage()
		return nil, fmt.Errorf("invalid dimensions")
	}

	if *concurrency < 1 {
		flag.Usage()
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	formatter := format.StandardFormat
	if *isCompact {
		formatter = format.CompactFormat
	}

	return &ParseResult{
		rows:        *rows,
		columns:     *columns,
		concurrency: *concurrency,
		formatter:   formatter,
	}, nil
}

func isDimensionValid(size int) bool {
	if size < 3 || size > 13 {
		return false
	}
	return true
}
