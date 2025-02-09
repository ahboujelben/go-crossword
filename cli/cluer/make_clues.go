package cluer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/ahboujelben/crossword/generator"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	System string `json:"system"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

var system = `
You are a crossword clue generator. You will generate a single concise
sentence for a given word. The crossword's difficulty should be normal.
`

func MakeClues(c *generator.Crossword) map[string]string {
	clues := make(map[string]string)
	mutex := sync.Mutex{}
	for word := generator.Word(c); word != nil; word = word.Next() {
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func(w string) {
			defer wg.Done()
			prompt := OllamaRequest{
				Model:  "llama3.1:8b",
				System: system,
				Prompt: w,
				Stream: false,
			}
			jsonData, err := json.Marshal(prompt)
			if err != nil {
				return
			}
			resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
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
			mutex.Lock()

			// remove quotes if present
			clue := result.Response
			if len(clue) > 1 && clue[0] == '"' && clue[len(clue)-1] == '"' {
				clue = clue[1 : len(clue)-1]
			}
			// remove trailing period
			if len(clue) > 1 && clue[len(clue)-1] == '.' {
				clue = clue[:len(clue)-1]
			}

			clues[w] = clue
			mutex.Unlock()
		}(string(word.GetValue()))
		wg.Wait()
	}
	return clues
}
