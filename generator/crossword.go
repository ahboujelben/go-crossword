package generator

import (
	"fmt"
	"math/rand"
)

type Crossword struct {
	Width  int
	Height int
	Data   []byte
}

func (c Crossword) IsFilled() bool {
	for _, letter := range c.Data {
		if letter == 0 {
			return false
		}
	}
	return true
}

func (c Crossword) Words() WordRefs {
	words := c.HorizontalWords()
	words = append(words, c.VerticalWords()...)
	return words
}

func (c Crossword) HorizontalWords() []WordRef {
	refs := make([]WordRef, 0)

	for y := 0; y < c.Height; y++ {
		word := []byte{}
		for x := 0; x < c.Width; x++ {
			if c.Data[y*c.Width+x] == '.' {
				if len(word) > 1 {
					refs = append(refs, NewWordRef(y*c.Width+x-len(word), len(word), Horizontal, &c))
				}
				word = []byte{}
			} else {
				word = append(word, c.Data[y*c.Width+x])
			}
		}
		if len(word) > 1 {
			refs = append(refs, NewWordRef(y*c.Width+c.Width-len(word), len(word), Horizontal, &c))
		}
	}

	return refs
}

func (c Crossword) VerticalWords() []WordRef {
	refs := make([]WordRef, 0)

	for x := 0; x < c.Width; x++ {
		word := []byte{}
		for y := 0; y < c.Height; y++ {
			if c.Data[y*c.Width+x] == '.' {
				if len(word) > 1 {
					refs = append(refs, NewWordRef((y-len(word))*c.Width+x, len(word), Vertical, &c))
				}
				word = []byte{}
			} else {
				word = append(word, c.Data[y*c.Width+x])
			}
		}
		if len(word) > 1 {
			refs = append(refs, NewWordRef((c.Height-len(word))*c.Width+x, len(word), Vertical, &c))
		}
	}

	return refs
}

func (c Crossword) Print() {
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			if c.Data[y*c.Width+x] == 0 {
				fmt.Print("*")
			} else {
				fmt.Printf("%c", c.Data[y*c.Width+x])
			}
		}
		fmt.Println()
	}
}

func (c Crossword) FillFromDict(wordDict WordDict) Crossword {
	words := c.Words().Sorted()

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
			for _, word := range c.Words() {
				if word.IsFilled() {
					if !wordDict.Contains(string(word.GetValue())) {
						valid = false
						break
					}
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
		fmt.Println()
	}
	return c
}

func NewCrosswordFromDict(width, height int, wordDict WordDict) Crossword {
	return NewEmptyCrossword(width, height).FillFromDict(wordDict)
}

func NewEmptyCrossword(width, height int) Crossword {
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
					data[i*width+j] = '.'
				}
			}
		} else {
			if i != 0 && i != height-1 {
				data[i*width+i] = '.'
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
				data[i*width+rand.Intn(width-offset[width]*2)+offset[width]] = '.'
			}
			if rand.Float64() < chance[width] {
				data[(rand.Intn(height-offset[width]*2)+offset[width])*width+i] = '.'
			}
		}
	}

	// replace any single letter words with empty space
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if data[y*width+x] == '.' {
				continue
			}
			if (x == 0 || data[y*width+x-1] == '.') &&
				(x == width-1 || data[y*width+x+1] == '.') &&
				(y == 0 || data[(y-1)*width+x] == '.') &&
				(y == height-1 || data[(y+1)*width+x] == '.') {
				data[y*width+x] = '.'
			}
		}
	}

	return Crossword{
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
		if content[i] != '.' && (content[i] < 'a' || content[i] > 'z') {
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
