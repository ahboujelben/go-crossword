package main

import "github.com/ahboujelben/crossword/generator"

func main() {
	wordDict := generator.NewWordDict("data/words.txt")
	crossword := generator.NewCrosswordFromDict(13, 13, wordDict)
	crossword.Print()
}
