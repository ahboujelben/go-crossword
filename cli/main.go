package main

import (
	"flag"
	"fmt"

	"github.com/ahboujelben/crossword/cli/format"
	"github.com/ahboujelben/crossword/generator"
)

func main() {
	config, err := parseArguments()
	if err != nil {
		return
	}

	wordDict := generator.NewWordDict("data/words.txt")
	crossword := generator.NewCrossword(config.columns, config.rows, wordDict)
	crossword.Print(config.formatter)
}

type config struct {
	columns   int
	rows      int
	formatter func(c *generator.Crossword)
}

func parseArguments() (*config, error) {
	columns := flag.Int("columns", 13, "number of columns in the crossword (valid values: [2, 13])")
	rows := flag.Int("rows", 13, "number of rows in the crossword (valid values: [2, 13])")
	isCompact := flag.Bool("compact", false, "prints each letter using one character")
	flag.Parse()

	if !isDimensionValid(*rows) || !isDimensionValid(*columns) {
		flag.Usage()
		return nil, fmt.Errorf("invalid dimensions")
	}

	formatter := format.StandardFormat
	if *isCompact {
		formatter = format.CompactFormat
	}

	return &config{
		columns:   *columns,
		rows:      *rows,
		formatter: formatter,
	}, nil
}

func isDimensionValid(size int) bool {
	if size < 2 || size > 13 {
		return false
	}
	return true
}
