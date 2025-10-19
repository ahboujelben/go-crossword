package main

import (
	"context"
	"log"

	"github.com/ahboujelben/go-crossword/cli/renderer"
	"github.com/ahboujelben/go-crossword/modules/crossword"
	"github.com/ahboujelben/go-crossword/modules/dictionary"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Input struct {
	Rows int `json:"rows" jsonschema:"the number of rows in the crossword"`
	Cols int `json:"cols" jsonschema:"the number of columns in the crossword"`
}

type Output struct {
	UnsolvedCrossword string `json:"unsolvedCrossword" jsonschema:"the crossword grid without the solution - to be printed as is"`
	SolvedCrossword   string `json:"solvedCrossword" jsonschema:"the crossword grid with the solution - to be printed as is"`
	RowWords          []Word `json:"rowWords" jsonschema:"the list of row words in the solved crossword"`
	ColumnWords       []Word `json:"columnWords" jsonschema:"the list of column words in the solved crossword"`
}

type Word struct {
	Value  string `json:"value"`
	Row    int    `json:"row"`
	Column int    `json:"column"`
}

func newRowWord(word *crossword.RowWordRef) Word {
	return Word{
		Value:  string(word.GetValue()),
		Row:    word.Row() + 1,
		Column: word.Column() + 1,
	}
}

func newColumnWord(word *crossword.ColumnWordRef) Word {
	return Word{
		Value:  string(word.GetValue()),
		Row:    word.Row() + 1,
		Column: word.Column() + 1,
	}
}

// isSizeValid checks if a crossword size is valid
func isSizeValid(size int) bool {
	return size >= 3 && size <= 15
}

func GenerateCrossword(ctx context.Context, req *mcp.CallToolRequest, input Input) (
	*mcp.CallToolResult,
	Output,
	error,
) {
	// Validate input dimensions
	if !isSizeValid(input.Rows) || !isSizeValid(input.Cols) {
		return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: "rows and cols must be between 3 and 15 inclusive",
					},
				},
				IsError: true,
			}, Output{
				UnsolvedCrossword: "",
				SolvedCrossword:   "",
				RowWords:          []Word{},
				ColumnWords:       []Word{},
			}, nil
	}

	result := crossword.NewCrossword(crossword.CrosswordConfig{
		Rows:     input.Rows,
		Cols:     input.Cols,
		Threads:  100,
		WordDict: dictionary.NewWordDictionary(),
	})

	c := result.Crossword

	unsolvedCrossword := renderer.NewStandardRenderer().RenderCrossword(c, false)
	solvedCrossword := renderer.NewStandardRenderer().RenderCrossword(c, true)

	rowWords := []Word{}
	for word := crossword.RowWord(c); word != nil; word = word.Next() {
		rowWords = append(rowWords, newRowWord(word))
	}

	columnWords := []Word{}
	for word := crossword.ColumnWord(c); word != nil; word = word.Next() {
		columnWords = append(columnWords, newColumnWord(word))
	}

	return nil,
		Output{
			UnsolvedCrossword: unsolvedCrossword,
			SolvedCrossword:   solvedCrossword,
			RowWords:          rowWords,
			ColumnWords:       columnWords,
		},
		nil
}

func main() {
	// Create a server with the get crossword tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "go-crossword", Version: "v1.0.0"}, nil)

	toolDescription := `
Generate a crossword puzzle with the specified dimensions.

IMPORTANT: When a user requests a crossword, always:
1. Display the unsolved crossword grid in monospace font and always print the newline characters between rows.
2. Generate an interesting clue for each word but without displaying the word.
3. Prefix each clue for across words with its position in the grid in this format: (Row: x, Col: y).
s. Prefix each clue for down words with its position in the grid in this format (Col: y, Row: x).
5. Only reveal the solved solution and the words when explicitly requested or if the user gives up.

The clues should be presented in a clear format with numbered clues for Across and Down.
`

	mcp.AddTool(server, &mcp.Tool{Name: "generate-crossword", Description: toolDescription}, GenerateCrossword)

	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
