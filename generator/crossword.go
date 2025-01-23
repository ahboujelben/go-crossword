package generator

import (
	"fmt"
	"math/rand"
	"sort"
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

func (c *Crossword) FillFromDict(wordDict WordDict) *Crossword {
	words := make([]WordRef, 0)
	for w := c.FirstWord(); w != nil; w = w.Next() {
		words = append(words, *w)
	}
	sort.Slice(words, func(i, j int) bool {
		return words[i].Length > words[j].Length
	})

	currentWordIndex := 0
	totalBacktracks := 0
	backtrackSteps := 3

	type StackElement struct {
		index int
		word  []byte
	}
	stack := []StackElement{}

	backtrack := func() {
		totalBacktracks++
		if totalBacktracks%10 == 0 {
			backtrackSteps += 3
		}
		for i := 0; i < backtrackSteps; i++ {
			prevWord := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			currentWordIndex = prevWord.index
			words[currentWordIndex].SetValue(prevWord.word)
			if len(stack) == 0 {
				backtrackSteps = 3
				totalBacktracks = 0
				break
			}
		}
	}

	for {
		if c.IsFilled() {
			break
		}

		currentWord := words[currentWordIndex]
		currentWordValue := currentWord.GetValue()

		// if the current word is already filled, move on to the next word
		// this can happen as a result of filling other words in the crossword
		if currentWord.IsFilled() {
			currentWordIndex++
			continue
		}

		// find all the fillers that are potential candidates for the current word
		candidates := wordDict.Candidates(currentWordValue)

		// if no candidates are found, backtrack
		if candidates.Len() == 0 {
			backtrack()
			continue
		}

		// Choose a random candidate that keeps the crossword valid
		found := false
		for candidates.Len() > 0 {
			// choose a random candidate
			e := candidates.Front()
			chosenIndex := rand.Intn(candidates.Len())
			for i := 0; i < chosenIndex; i++ {
				e = e.Next()
			}
			candidate := e.Value.(string)

			// check if the crossword is still valid once the current word is
			// filled with the candidate
			currentWord.SetValue([]byte(candidate))
			valid := true
			for word := c.FirstWord(); word != nil; word = word.Next() {
				if word.IsFilled() && !wordDict.Contains(string(word.GetValue())) {
					valid = false
					break
				}
			}

			if valid {
				found = true
				break
			}
		}

		if !found {
			currentWord.SetValue(currentWordValue)
			backtrack()
			continue
		}

		stack = append(stack, StackElement{index: currentWordIndex, word: currentWordValue})
		currentWordIndex++

		fmt.Print("\033[H\033[2J")
		c.Print()
	}
	return c
}

func NewCrosswordFromDict(width, height int, wordDict WordDict) *Crossword {
	return NewEmptyCrossword(width, height).FillFromDict(wordDict)
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
		if i%2 == 1 {
			for j := 0; j < width; j++ {
				if j%2 == 1 {
					data[i*width+j] = blankLetter
				}
			}
		} else {
			if i != 0 && i != height-1 {
				data[i*width+i] = blankLetter
			}
		}
		if width > 7 && (i == 0 || i == height-1) {
			chance := map[int]float64{
				9:  0.3,
				11: 0.6,
				13: 0.9,
				15: 1,
			}
			offset := map[int]int{
				9:  3,
				11: 3,
				13: 3,
				15: 5,
			}
			if rand.Float64() < chance[width] {
				data[i*width+rand.Intn(width-offset[width]*2)+offset[width]] = blankLetter
			}
			if rand.Float64() < chance[width] {
				data[(rand.Intn(height-offset[width]*2)+offset[width])*width+i] = blankLetter
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
		panic(fmt.Sprintf("expected length %d but got %d", length, len(content)))
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
