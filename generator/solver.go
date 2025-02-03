package generator

import (
	"context"
	"math/rand"
	"slices"
	"sort"
)

// starting with an empty crossword, try to fill the crossword word by word,
// starting with the longest ones. if stuck or we ended up creating
// non-existent words, backtrack and try again.
func generateCrossword(ctx context.Context, rows, columns int, wordDict WordDict, solvedCrossword chan *Crossword) {
	crossword := NewEmptyCrossword(rows, columns)
	crawler := newCrosswordCrawler(crossword)

	for {
		// abort if the context is cancelled - a solution has already been found
		select {
		case <-ctx.Done():
			return
		default:
		}

		// if the whole crossword is filled then a solution has been found
		if crossword.isFilled() {
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

		// find possible candidates for the current word based on the current
		// state of the crossword
		candidates := wordDict.Candidates(currentWordValue)
		// exclude words that are already in the crossword
		candidates = slices.DeleteFunc(candidates, func(e int) bool {
			_, exists := crawler.wordsSoFar[wordDict.allWords[e]]
			return exists
		})

		if len(candidates) == 0 {
			crawler.backtrack()
			continue
		}

		candidate := wordDict.allWords[candidates[rand.Intn(len(candidates))]]
		crawler.pushToStack(currentWordValue)
		currentWord.SetValue([]byte(candidate))
		crawler.storeWord(candidate)
		crawler.goToNextWord()
	}
}

type crosswordCrawler struct {
	words            []WordRef
	stack            []wordStack
	wordsSoFar       map[string]struct{}
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
		return words[i].length > words[j].length
	})
	return &crosswordCrawler{
		words:            words,
		stack:            []wordStack{},
		wordsSoFar:       make(map[string]struct{}),
		currentWordIndex: 0,
		totalBacktracks:  0,
		backtrackSteps:   3,
	}
}

func (c *crosswordCrawler) pushToStack(value []byte) {
	c.stack = append(c.stack, wordStack{index: c.currentWordIndex, word: value})
}

func (c *crosswordCrawler) storeWord(value string) {
	c.wordsSoFar[value] = struct{}{}
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
		wordToBeDeleted := string(c.words[prevWord.index].GetValue())
		delete(c.wordsSoFar, wordToBeDeleted)
		c.currentWordIndex = prevWord.index
		c.words[c.currentWordIndex].SetValue(prevWord.word)
		if len(c.stack) == 0 {
			c.backtrackSteps = 3
			c.totalBacktracks = 0
			break
		}
	}
}
