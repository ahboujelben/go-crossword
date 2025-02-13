package cluer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"

	"github.com/ahboujelben/go-crossword/generator"
)

type ollamaRequest struct {
	Model   string               `json:"model"`
	System  string               `json:"system"`
	Prompt  string               `json:"prompt"`
	Stream  bool                 `json:"stream"`
	Options ollamaRequestOptions `json:"options"`
}

type ollamaRequestOptions struct {
	Seed int64 `json:"seed"`
}

type CluesResult struct {
	Clues map[string]string
	Seed  int64
}

func newCluesResult(clues map[string]string, seed int64) CluesResult {
	return CluesResult{
		Clues: clues,
		Seed:  seed,
	}
}

type ClueDifficulty string

const (
	ClueDifficultyNormal  ClueDifficulty = "normal"
	ClueDifficultyCryptic ClueDifficulty = "cryptic"
)

func getSystemPrompt(difficulty ClueDifficulty) string {
	return fmt.Sprintf(`
You are a crossword clue generator. You will generate a single concise
sentence for a given word. The crossword's difficulty should be %s.
`, difficulty)
}

type CluesConfig struct {
	Seed        int64
	Difficulty  ClueDifficulty
	OllamaModel string
	OllamaUrl   string
}

func NewClues(c *generator.Crossword, config CluesConfig) CluesResult {
	if config.Seed == 0 {
		config.Seed = rand.Int63()
	}

	clues := make(map[string]string)
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}

	for word := generator.Word(c); word != nil; word = word.Next() {
		wg.Add(1)
		go func(w string) {
			defer wg.Done()
			prompt := ollamaRequest{
				Model:  config.OllamaModel,
				System: getSystemPrompt(config.Difficulty),
				Prompt: fmt.Sprintf("the word is: `%s`", w),
				Stream: false,
				Options: ollamaRequestOptions{
					Seed: config.Seed,
				},
			}
			jsonData, err := json.Marshal(prompt)
			if err != nil {
				return
			}
			resp, err := http.Post(fmt.Sprintf("%s/api/generate", config.OllamaUrl), "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				return
			}
			defer resp.Body.Close()
			var result struct {
				Response string `json:"response"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return
			}

			// remove quotes if present
			clue := result.Response
			if len(clue) > 1 && clue[0] == '"' && clue[len(clue)-1] == '"' {
				clue = clue[1 : len(clue)-1]
			}
			// remove trailing period
			if len(clue) > 1 && clue[len(clue)-1] == '.' {
				clue = clue[:len(clue)-1]
			}

			mutex.Lock()
			clues[w] = clue
			mutex.Unlock()
		}(string(word.GetValue()))
	}

	wg.Wait()
	return newCluesResult(clues, config.Seed)
}
