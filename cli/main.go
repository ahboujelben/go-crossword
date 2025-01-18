package main

import (
	"github.com/ahboujelben/crosswords/generator"
)

func main() {
	wordDict := generator.NewWordDict("data/words.txt")
	generator.NewCrosswordFromDict(13, 13, wordDict)
	// crossword := generator.NewCrossword().FromString(5, 5, "hi.oh...ah..uh..wee.love.")
	// crossword.Print()
	// for _, word := range crossword.Words() {
	// 	println(string(word))
	// }

}
