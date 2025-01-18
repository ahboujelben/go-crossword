package generator

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
)

type WordDict struct {
	AllWords      map[string]struct{}
	LengthMapping map[int][]string
	LetterMapping map[WordDictKey]map[string]struct{}
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

	// iterate over the words from data/words.txt and build the dictionary
	dict := WordDict{
		AllWords:      map[string]struct{}{},
		LengthMapping: map[int][]string{},
		LetterMapping: map[WordDictKey]map[string]struct{}{},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		dict.AllWords[word] = struct{}{}
		dict.LengthMapping[len(word)] = append(dict.LengthMapping[len(word)], word)
		for i := 0; i < len(word); i++ {
			key := WordDictKey{letter: word[i], pos: i}
			if dict.LetterMapping[key] == nil {
				dict.LetterMapping[key] = map[string]struct{}{}
			}
			dict.LetterMapping[key][word] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("error reading words file: %s", err))
	}

	return dict

}

func (wd WordDict) Contains(word string) bool {
	_, exists := wd.AllWords[word]
	return exists
}

func (wd WordDict) Candidates(word []byte) *list.List {
	wordSets := list.New()
	for _, word := range wd.LengthMapping[len(word)] {
		wordSets.PushBack(word)
	}

	for i, letter := range word {
		if letter != 0 {
			key := WordDictKey{letter: letter, pos: i}
			currentSet := wd.LetterMapping[key]
			for e := wordSets.Front(); e != nil; {
				next := e.Next()
				if _, exists := currentSet[e.Value.(string)]; !exists {
					wordSets.Remove(e)
				}
				e = next
			}
		}
	}
	return wordSets
}
