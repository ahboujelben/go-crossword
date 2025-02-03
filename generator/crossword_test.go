package generator_test

import (
	"fmt"
	"testing"

	"github.com/ahboujelben/crossword/generator"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCrossword(t *testing.T) {
	wordDict := generator.NewWordDict("../data/words.txt")
	for rows := 3; rows <= 13; rows++ {
		for columns := 3; columns <= 13; columns++ {
			c := columns
			r := rows
			t.Run(fmt.Sprintf("Rows=%d_Columns=%d", rows, columns), func(t *testing.T) {
				t.Parallel()
				crossword := generator.NewCrossword(generator.CrosswordConfig{
					Rows:        r,
					Columns:     c,
					Concurrency: 100,
					WordDict:    wordDict,
				})

				assert.True(t, crossword.IsFilled())
				for word := generator.ColumnWord(crossword); word != nil; word = word.Next() {
					wordValue := string(word.GetValue())
					assert.True(t, wordDict.Contains(wordValue))
				}
			})
		}
	}
}
