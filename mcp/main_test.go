package main

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGenerateCrossword(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	t.Run("valid input generates a crossword", func(t *testing.T) {
		input := Input{Rows: 5, Cols: 5}
		result, output, err := GenerateCrossword(ctx, req, input)

		if err != nil {
			t.Fatalf("GenerateCrossword() returned an unexpected error: %v", err)
		}

		if result != nil {
			t.Fatalf("Expected a nil mcp.CallToolResult for valid input, but got: %+v", result)
		}

		if output.UnsolvedCrossword == "" {
			t.Error("UnsolvedCrossword should not be empty for valid input")
		}

		if output.SolvedCrossword == "" {
			t.Error("SolvedCrossword should not be empty for valid input")
		}

		if len(output.RowWords) == 0 {
			t.Error("RowWords should not be empty for a valid crossword")
		}

		if len(output.ColumnWords) == 0 {
			t.Error("ColumnWords should not be empty for a valid crossword")
		}
	})

	t.Run("invalid input returns an error result", func(t *testing.T) {
		testCases := []struct {
			name  string
			input Input
		}{
			{"rows too small", Input{Rows: 2, Cols: 5}},
			{"cols too small", Input{Rows: 5, Cols: 1}},
			{"rows too large", Input{Rows: 16, Cols: 5}},
			{"cols too large", Input{Rows: 5, Cols: 20}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				result, output, err := GenerateCrossword(ctx, req, tc.input)

				if err != nil {
					t.Fatalf("GenerateCrossword() returned an unexpected error: %v", err)
				}

				if result == nil {
					t.Fatal("Expected a non-nil mcp.CallToolResult for invalid input, but got nil")
				}

				if !result.IsError {
					t.Error("Expected result.IsError to be true for invalid input")
				}

				expectedMsg := "rows and cols must be between 3 and 15 inclusive"
				if len(result.Content) == 0 {
					t.Fatal("Expected result.Content to have at least one item")
				}
				textContent, ok := result.Content[0].(*mcp.TextContent)
				if !ok {
					t.Fatal("Expected result.Content[0] to be of type *mcp.TextContent")
				}

				if textContent.Text != expectedMsg {
					t.Errorf("Expected error message '%s', but got '%s'", expectedMsg, textContent.Text)
				}

				if output.UnsolvedCrossword != "" || output.SolvedCrossword != "" || len(output.RowWords) > 0 || len(output.ColumnWords) > 0 {
					t.Error("Expected an empty Output struct for invalid input")
				}
			})
		}
	})
}
