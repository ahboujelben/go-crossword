package generator

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
)

const Blank = '.'

type Crossword struct {
	columns int
	rows    int
	data    []byte
}

type CrosswordConfig struct {
	Columns     int
	Rows        int
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
			generateCrossword(ctx, config.Columns, config.Rows, config.WordDict, solvedCrossword)
		}()
	}

	wg.Wait()

	return <-solvedCrossword
}

func NewEmptyCrossword(columns, rows int) *Crossword {
	if columns < 1 {
		panic(fmt.Sprintf("invalid columns: %d", columns))
	}
	if rows < 1 {
		panic(fmt.Sprintf("invalid rows: %d", rows))
	}

	data := make([]byte, columns*rows)

	// create blank squares based on specific conditions
	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			if i%2 == 1 && (i+j)%2 == 0 {
				data[i*columns+j] = Blank
			}
			if i == 0 && j%2 == 0 && rand.Float64() < 0.75 {
				data[j+rand.Intn(rows)*columns] = Blank
			}
		}

		if i%2 == 0 && columns > 7 {
			if rand.Float64() < 0.75 {
				data[i*columns+rand.Intn(columns)] = Blank
			}
		}
	}

	// replace any single letter words with empty space
	for y := 0; y < rows; y++ {
		for x := 0; x < columns; x++ {
			if data[y*columns+x] == Blank {
				continue
			}
			if (x == 0 || data[y*columns+x-1] == Blank) &&
				(x == columns-1 || data[y*columns+x+1] == Blank) &&
				(y == 0 || data[(y-1)*columns+x] == Blank) &&
				(y == rows-1 || data[(y+1)*columns+x] == Blank) {
				data[y*columns+x] = Blank
			}
		}
	}

	return &Crossword{
		columns: columns,
		rows:    rows,
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

func (c *Crossword) isFilled() bool {
	for letter := CrosswordLetter(c); letter != nil; letter = letter.Next() {
		if letter.IsEmpty() {
			return false
		}
	}
	return true
}
