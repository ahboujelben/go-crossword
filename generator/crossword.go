package generator

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"sync"
)

const blankLetter = '.'

type Crossword struct {
	Width  int
	Height int
	Data   []byte
}

func (c *Crossword) FirstLetter() *CrosswordLetterRef {
	return NewCrosswordLetterRef(0, c)
}

func (c *Crossword) FirstWord() *WordRef {
	return NewHorizontalWordRef(0, c)
}

func (c *Crossword) IsFilled() bool {
	for letter := c.FirstLetter(); letter != nil; letter = letter.Next() {
		if !letter.IsFilled() {
			return false
		}
	}
	return true
}

func (c *Crossword) Print() {
	fmt.Println()
	for letter := c.FirstLetter(); letter != nil; letter = letter.Next() {
		if letter.GetValue() == blankLetter {
			fmt.Print("# ")
		} else if letter.IsFilled() {
			fmt.Printf("%c ", letter.GetValue())
		} else {
			fmt.Print("* ")
		}
		if letter.Pos%c.Width == c.Width-1 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func NewCrosswordFromDict(width, height int, wordDict WordDict) *Crossword {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	filledCrossword := make(chan *Crossword, 1)

	numThreads := 100
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		// threadNumber := i
		go func() {
			defer wg.Done()
			defer cancel()
			// fmt.Printf("Thread %d started", threadNumber)
			// fmt.Println()
			crossword := NewEmptyCrossword(width, height)
			crossword.solve(ctx, wordDict, filledCrossword)
			// fmt.Printf("Thread %d ended", threadNumber)
			// fmt.Println()
		}()
	}

	wg.Wait()

	return <-filledCrossword
}

// try to fill the crossword word by word, starting with the longest ones.
// if stuck or we ended up creating non-existent words, backtrack and try again.
func (c *Crossword) solve(ctx context.Context, wordDict WordDict, filledCrossword chan *Crossword) {
	words := make([]WordRef, 0)
	for w := c.FirstWord(); w != nil; w = w.Next() {
		words = append(words, *w)
	}
	sort.Slice(words, func(i, j int) bool {
		return words[i].Length > words[j].Length
	})

	crawler := newCrosswordCrawler(words)

	for {
		// abort if the context is cancelled - a solution has already been found
		select {
		case <-ctx.Done():
			return
		default:
		}

		// if the whole crossword is filled then a solution has been found
		if c.IsFilled() {
			filledCrossword <- c
			return
		}

		currentWord := words[crawler.currentWordIndex]
		currentWordValue := currentWord.GetValue()

		if currentWord.IsFilled() {
			if !wordDict.Contains(string(currentWordValue)) {
				crawler.backtrack()
				continue
			}
			crawler.currentWordIndex++
			continue
		}

		candidates := wordDict.Candidates(currentWordValue)
		if len(candidates) == 0 {
			crawler.backtrack()
			continue
		}

		crawler.stack = append(crawler.stack, wordStack{index: crawler.currentWordIndex, word: currentWordValue})
		currentWord.SetValue([]byte(candidates[rand.Intn(len(candidates))]))
		crawler.currentWordIndex++
	}
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
				data[i*width+j] = blankLetter
			}
		}

		if i%2 == 0 && width > 7 {
			if rand.Float64() < 0.75 {
				data[i*width+rand.Intn(width)] = blankLetter
			}
			if rand.Float64() < 0.75 {
				data[i+rand.Intn(height)*width] = blankLetter
			}
		}
	}

	// replace any single letter words with empty space
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if data[y*width+x] == blankLetter {
				continue
			}
			if (x == 0 || data[y*width+x-1] == blankLetter) &&
				(x == width-1 || data[y*width+x+1] == blankLetter) &&
				(y == 0 || data[(y-1)*width+x] == blankLetter) &&
				(y == height-1 || data[(y+1)*width+x] == blankLetter) {
				data[y*width+x] = blankLetter
			}
		}
	}

	return &Crossword{
		Width:  width,
		Height: height,
		Data:   data,
	}
}

func NewCrosswordFromString(width, height int, content string) Crossword {
	if width < 1 {
		panic(fmt.Sprintf("invalid width: %d", width))
	}

	if height < 1 {
		panic(fmt.Sprintf("invalid height: %d", height))
	}

	length := width * height
	if len(content) != length {
		panic(fmt.Sprintf("unexpected content length (expected=%d,found=%d)", length, len(content)))
	}

	data := make([]byte, len(content))
	for i := range content {
		if content[i] != blankLetter && (content[i] < 'a' || content[i] > 'z') {
			panic(fmt.Sprintf("invalid character at position %d: %c", i, content[i]))
		}
		data[i] = content[i]
	}

	return Crossword{
		Width:  width,
		Height: height,
		Data:   data,
	}
}
