package generator

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type WordDict struct {
	allWords  []string
	wordSet   map[string]struct{}
	lengthMap map[int][]int
	letterMap map[WordDictKey]map[int]struct{}
}

type WordDictKey struct {
	letter byte
	pos    int
}

func NewWordDict(dictPath string) WordDict {
	file, err := os.Open(dictPath)
	if err != nil {
		panic(fmt.Sprintf("failed to open words file: %s", err))
	}
	defer file.Close()

	dict := WordDict{
		allWords:  []string{},
		wordSet:   map[string]struct{}{},
		lengthMap: map[int][]int{},
		letterMap: map[WordDictKey]map[int]struct{}{},
	}

	scanner := bufio.NewScanner(file)
	wordIndex := 0
	for scanner.Scan() {
		word := scanner.Text()
		dict.allWords = append(dict.allWords, word)
		dict.wordSet[word] = struct{}{}
		dict.lengthMap[len(word)] = append(dict.lengthMap[len(word)], wordIndex)
		for i := 0; i < len(word); i++ {
			key := WordDictKey{letter: word[i], pos: i}
			if _, exists := dict.letterMap[key]; !exists {
				dict.letterMap[key] = map[int]struct{}{}
			}
			dict.letterMap[key][wordIndex] = struct{}{}
		}
		wordIndex++
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("error reading dictionary file: %s", err))
	}

	return dict

}

func (wd WordDict) Contains(word string) bool {
	_, exists := wd.wordSet[word]
	return exists
}

func (wd WordDict) Candidates(word []byte) []int {
	candidates := make([]int, len(wd.lengthMap[len(word)]))
	copy(candidates, wd.lengthMap[len(word)])

	for i, letter := range word {
		if letter != 0 {
			key := WordDictKey{letter: letter, pos: i}
			currentSet := wd.letterMap[key]
			candidates = slices.DeleteFunc(candidates, func(e int) bool {
				_, exists := currentSet[e]
				return !exists
			})
		}
	}
	return candidates
}
