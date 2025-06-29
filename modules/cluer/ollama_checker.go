package cluer

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ollamaTags struct {
	Models []ollamaModel `json:"models"`
}

type ollamaModel struct {
	Name string `json:"name"`
}

func CheckOllamaServer(ollamaUrl string, ollamaModel string) error {
	// query {ollamaUrl}/api/tags which returns a json of type ollamaTags
	// throw an error if the server is not reachable or the model is not found
	resp, err := http.Get(fmt.Sprintf("%s/api/tags", ollamaUrl))
	if err != nil {
		return fmt.Errorf("could not reach ollama server at %s", ollamaUrl)
	}
	defer resp.Body.Close()
	var tags ollamaTags
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return fmt.Errorf("could not decode ollama tags response")
	}
	for _, model := range tags.Models {
		if model.Name == ollamaModel {
			return nil
		}
	}
	return fmt.Errorf("Ollama does not provide the model `%s`", ollamaModel)
}
