package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ahboujelben/crossword/cli/format"
	"github.com/ahboujelben/crossword/generator"
)

func main() {
	parseResult, err := parseArguments()
	if err != nil {
		return
	}

	crossword := generator.NewCrossword(parseResult.config)
	crossword.Print(parseResult.formatter)
}

type ParseResult struct {
	config    generator.CrosswordConfig
	formatter func(c *generator.Crossword)
}

func parseArguments() (*ParseResult, error) {
	rows := flag.Int("rows", 13, "number of rows in the crossword (valid values: [3, 13])")
	columns := flag.Int("cols", 13, "number of columns in the crossword (valid values: [3, 13])")
	concurrency := flag.Int("conc", 100, "number of goroutines to use (valid values: >= 1)")
	dictionaryPath := flag.String("dict", "data/words.txt", "path to the dictionary file with the words to be used to fill the crossword")
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

	if _, err := os.Stat(*dictionaryPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("dictionary file not found: %s", *dictionaryPath)
	}

	wordDict := generator.NewWordDict(*dictionaryPath)

	formatter := format.StandardFormat
	if *isCompact {
		formatter = format.CompactFormat
	}

	return &ParseResult{
		config: generator.CrosswordConfig{
			Rows:        *rows,
			Columns:     *columns,
			Concurrency: *concurrency,
			WordDict:    wordDict,
		},
		formatter: formatter,
	}, nil
}

func isDimensionValid(size int) bool {
	if size < 3 || size > 13 {
		return false
	}
	return true
}
