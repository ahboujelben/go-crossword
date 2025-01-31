package generator

type wordStack struct {
	index int
	word  []byte
}
type crosswordCrawler struct {
	words            []WordRef
	stack            []wordStack
	currentWordIndex int
	totalBacktracks  int
	backtrackSteps   int
}

func newCrosswordCrawler(words []WordRef) *crosswordCrawler {
	return &crosswordCrawler{
		words:            words,
		stack:            []wordStack{},
		currentWordIndex: 0,
		totalBacktracks:  0,
		backtrackSteps:   3,
	}
}

func (e *crosswordCrawler) backtrack() {
	e.totalBacktracks++
	if e.totalBacktracks%10 == 0 {
		e.backtrackSteps += 3
	}
	for i := 0; i < e.backtrackSteps; i++ {
		prevWord := e.stack[len(e.stack)-1]
		e.stack = e.stack[:len(e.stack)-1]
		e.currentWordIndex = prevWord.index
		e.words[e.currentWordIndex].SetValue(prevWord.word)
		if len(e.stack) == 0 {
			e.backtrackSteps = 3
			e.totalBacktracks = 0
			break
		}
	}
}
