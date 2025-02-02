package generator

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
)

const blank = '.'

type Crossword struct {
	Width  int
	Height int
	data   []byte
}

func NewCrossword(width, height int, wordDict WordDict) *Crossword {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	solvedCrossword := make(chan *Crossword, 1)

	// Generating a random crossword can take an unpredictable amount of time,
	// depending on the initial crossword configuration and the words that are
	// tried. To speed up the process, we run multiple goroutines to generate
	// crosswords and return the first one that is solved. This typically takes
	// less than a second to generate a 13x13 crossword.
	numThreads := 100
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			generateCrossword(ctx, width, height, wordDict, solvedCrossword)
		}()
	}

	wg.Wait()

	return <-solvedCrossword
}

func NewEmptyCrossword(width, height int) *Crossword {
	if width < 1 {
		panic(fmt.Sprintf("invalid width: %d", width))
	}
	if height < 1 {
		panic(fmt.Sprintf("invalid height: %d", height))
	}

	data := make([]byte, width*height)

	// create blank squares based on specific conditions
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if i%2 == 1 && (i+j)%2 == 0 {
				data[i*width+j] = blank
			}
		}

		if i%2 == 0 && width > 7 {
			if rand.Float64() < 0.75 {
				data[i*width+rand.Intn(width)] = blank
			}
			if rand.Float64() < 0.75 {
				data[i+rand.Intn(height)*width] = blank
			}
		}
	}

	// replace any single letter words with empty space
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if data[y*width+x] == blank {
				continue
			}
			if (x == 0 || data[y*width+x-1] == blank) &&
				(x == width-1 || data[y*width+x+1] == blank) &&
				(y == 0 || data[(y-1)*width+x] == blank) &&
				(y == height-1 || data[(y+1)*width+x] == blank) {
				data[y*width+x] = blank
			}
		}
	}

	return &Crossword{
		Width:  width,
		Height: height,
		data:   data,
	}
}

func (c *Crossword) Print() {
	fmt.Println()
	for letter := CrosswordLetter(c); letter != nil; letter = letter.Next() {
		if letter.GetValue() == blank {
			fmt.Print("# ")
		} else if letter.IsEmpty() {
			fmt.Print("* ")
		} else {
			fmt.Printf("%c ", letter.GetValue())
		}
		if letter.Pos%c.Width == c.Width-1 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func (c *Crossword) IsFilled() bool {
	for letter := CrosswordLetter(c); letter != nil; letter = letter.Next() {
		if letter.IsEmpty() {
			return false
		}
	}
	return true
}
