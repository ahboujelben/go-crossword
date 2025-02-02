package generator

import (
	"context"
	"math/rand"
	"sort"
)

// starting with an empty crossword, try to fill the crossword word by word,
// starting with the longest ones. if stuck or we ended up creating
// non-existent words, backtrack and try again.
func generateCrossword(ctx context.Context, width, height int, wordDict WordDict, solvedCrossword chan *Crossword) {
	crossword := NewEmptyCrossword(width, height)
	crawler := newCrosswordCrawler(crossword)

	for {
		// abort if the context is cancelled - a solution has already been found
		select {
		case <-ctx.Done():
			return
		default:
		}

		// if the whole crossword is filled then a solution has been found
		if crossword.IsFilled() {
			solvedCrossword <- crossword
			return
		}

		currentWord := crawler.currentWord()
		currentWordValue := currentWord.GetValue()

		if currentWord.IsFilled() {
			if !wordDict.Contains(string(currentWordValue)) {
				crawler.backtrack()
				continue
			}
			crawler.goToNextWord()
			continue
		}

		candidates := wordDict.Candidates(currentWordValue)
		if len(candidates) == 0 {
			crawler.backtrack()
			continue
		}

		crawler.pushToStack(currentWordValue)
		currentWord.SetValue([]byte(wordDict.allWords[candidates[rand.Intn(len(candidates))]]))
		crawler.goToNextWord()
	}
}

type crosswordCrawler struct {
	words            []WordRef
	stack            []wordStack
	currentWordIndex int
	totalBacktracks  int
	backtrackSteps   int
}

type wordStack struct {
	index int
	word  []byte
}

func newCrosswordCrawler(c *Crossword) *crosswordCrawler {
	words := make([]WordRef, 0)
	for w := Word(c); w != nil; w = w.Next() {
		words = append(words, *w)
	}
	sort.Slice(words, func(i, j int) bool {
		return words[i].Length > words[j].Length
	})
	return &crosswordCrawler{
		words:            words,
		stack:            []wordStack{},
		currentWordIndex: 0,
		totalBacktracks:  0,
		backtrackSteps:   3,
	}
}

func (c *crosswordCrawler) pushToStack(value []byte) {
	c.stack = append(c.stack, wordStack{index: c.currentWordIndex, word: value})
}

func (c *crosswordCrawler) currentWord() *WordRef {
	return &c.words[c.currentWordIndex]
}

func (c *crosswordCrawler) goToNextWord() {
	c.currentWordIndex++
}

func (c *crosswordCrawler) backtrack() {
	c.totalBacktracks++
	if c.totalBacktracks%10 == 0 {
		c.backtrackSteps += 3
	}
	for i := 0; i < c.backtrackSteps; i++ {
		prevWord := c.stack[len(c.stack)-1]
		c.stack = c.stack[:len(c.stack)-1]
		c.currentWordIndex = prevWord.index
		c.words[c.currentWordIndex].SetValue(prevWord.word)
		if len(c.stack) == 0 {
			c.backtrackSteps = 3
			c.totalBacktracks = 0
			break
		}
	}
}
