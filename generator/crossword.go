package generator

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
)

const blank = '.'

type Crossword struct {
	rows    int
	columns int
	data    []byte
}

type CrosswordConfig struct {
	Rows        int
	Columns     int
	Concurrency int
	WordDict    WordDict
}

func NewCrossword(config CrosswordConfig) *Crossword {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	solvedCrossword := make(chan *Crossword, 1)

	// Generating a random crossword can take an unpredictable amount of time,
	// depending on the initial crossword configuration and the words that are
	// tried. To speed up the process, we run multiple goroutines to generate
	// crosswords and return the first one that is solved. This typically takes
	// less than a second to generate a 13x13 crossword.
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			generateCrossword(ctx, config.Rows, config.Columns, config.WordDict, solvedCrossword)
		}()
	}

	wg.Wait()

	return <-solvedCrossword
}

func NewEmptyCrossword(rows, columns int) *Crossword {
	if rows < 1 {
		panic(fmt.Sprintf("invalid rows: %d", rows))
	}
	if columns < 1 {
		panic(fmt.Sprintf("invalid columns: %d", columns))
	}

	data := make([]byte, columns*rows)

	// create blank squares based on specific conditions
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if i%2 == 1 && (i+j)%2 == 0 {
				data[i*columns+j] = blank
			}
			if i == 0 && j%2 == 0 && rand.Float64() < 0.75 {
				data[j+rand.Intn(rows)*columns] = blank
			}
		}

		if i%2 == 0 && columns > 7 {
			if rand.Float64() < 0.75 {
				data[i*columns+rand.Intn(columns)] = blank
			}
		}
	}

	// replace any single letter words with empty space
	for y := 0; y < rows; y++ {
		for x := 0; x < columns; x++ {
			if data[y*columns+x] == blank {
				continue
			}
			if (x == 0 || data[y*columns+x-1] == blank) &&
				(x == columns-1 || data[y*columns+x+1] == blank) &&
				(y == 0 || data[(y-1)*columns+x] == blank) &&
				(y == rows-1 || data[(y+1)*columns+x] == blank) {
				data[y*columns+x] = blank
			}
		}
	}

	return &Crossword{
		rows:    rows,
		columns: columns,
		data:    data,
	}
}

func (c *Crossword) Print(formatter func(c *Crossword)) {
	formatter(c)
}

func (c *Crossword) Columns() int {
	return c.columns
}

func (c *Crossword) Rows() int {
	return c.rows
}

func (c *Crossword) IsFilled() bool {
	for letter := CrosswordLetter(c); letter != nil; letter = letter.Next() {
		if letter.IsEmpty() {
			return false
		}
	}
	return true
}
