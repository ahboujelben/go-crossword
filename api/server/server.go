package server

import (
	"fmt"
	"strconv"

	"github.com/ahboujelben/go-crossword/modules/cluer"
	"github.com/ahboujelben/go-crossword/modules/crossword"
	"github.com/ahboujelben/go-crossword/modules/dictionary"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// ServerConfig represents the configuration for the API server
type ServerConfig struct {
	Port        int
	OllamaURL   string
	OllamaModel string
}

// Server represents the API server
type Server struct {
	app    *fiber.App
	config ServerConfig
}

// CrosswordResponse is the JSON response structure for the crossword endpoint
type CrosswordResponse struct {
	CrosswordSeed int64           `json:"crosswordSeed"`
	CluesSeed     int64           `json:"cluesSeed,omitempty"`
	Rows          int             `json:"rows"`
	Cols          int             `json:"cols"`
	Data          string          `json:"data"`
	Words         []CrosswordWord `json:"words"`
}

// CrosswordWord represents a word in the crossword
type CrosswordWord struct {
	Pos       int    `json:"pos"`
	Direction string `json:"direction"`
	Value     string `json:"value,omitempty"`
	Clue      string `json:"clue,omitempty"`
}

// NewServer creates a new API server with the given configuration
func NewServer(config ServerConfig) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New())

	return &Server{
		app:    app,
		config: config,
	}
}

// Start starts the API server
func (s *Server) Start() error {
	s.setupRoutes()
	return s.app.Listen(fmt.Sprintf(":%d", s.config.Port))
}

// setupRoutes configures the API routes
func (s *Server) setupRoutes() {
	s.app.Get("/api/crossword", s.handleGetCrossword)
}

// handleGetCrossword handles the GET /api/crossword endpoint
func (s *Server) handleGetCrossword(c *fiber.Ctx) error {
	// Parse and validate query parameters
	rows, cols, seed, withClues, cluesSeed, unsolved, cryptic, err := parseAndValidateParams(c)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Generate the crossword
	crosswordResult := crossword.NewCrossword(crossword.CrosswordConfig{
		Rows:     rows,
		Cols:     cols,
		Seed:     seed,
		Threads:  100, // Default value as used in main.go
		WordDict: dictionary.NewWordDictionary(),
	})

	// Generate clues if requested
	var cluesResult *cluer.CluesResult
	if withClues {
		if err := cluer.CheckOllamaServer(s.config.OllamaURL, s.config.OllamaModel); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Ollama server error: %v", err))
		}

		clues := cluer.NewClues(crosswordResult.Crossword, cluer.CluesConfig{
			Seed:        cluesSeed,
			Cryptic:     cryptic,
			OllamaModel: s.config.OllamaModel,
			OllamaUrl:   s.config.OllamaURL,
		})
		cluesResult = &clues
	}

	// Prepare the response
	response := prepareCrosswordResponse(crosswordResult, cluesResult, unsolved)

	return c.JSON(response)
}

// parseAndValidateParams parses and validates the query parameters
func parseAndValidateParams(c *fiber.Ctx) (rows, cols int, seed int64, withClues bool, cluesSeed int64, unsolved bool, cryptic bool, err error) {
	// Parse rows parameter
	rowsStr := c.Query("rows", "13") // Default to 13 if not provided
	rows, err = strconv.Atoi(rowsStr)
	if err != nil {
		return 0, 0, 0, false, 0, false, false, fmt.Errorf("invalid rows parameter: %w", err)
	}
	if !isSizeValid(rows) {
		return 0, 0, 0, false, 0, false, false, fmt.Errorf("rows must be between 3 and 15")
	}

	// Parse cols parameter
	colsStr := c.Query("cols", "13") // Default to 13 if not provided
	cols, err = strconv.Atoi(colsStr)
	if err != nil {
		return 0, 0, 0, false, 0, false, false, fmt.Errorf("invalid cols parameter: %w", err)
	}
	if !isSizeValid(cols) {
		return 0, 0, 0, false, 0, false, false, fmt.Errorf("cols must be between 3 and 15")
	}

	// Parse seed parameter
	seedStr := c.Query("seed", "0") // Default to 0 (random seed) if not provided
	seed, err = strconv.ParseInt(seedStr, 10, 64)
	if err != nil {
		return 0, 0, 0, false, 0, false, false, fmt.Errorf("invalid seed parameter: %w", err)
	}
	if !isSeedValid(seed) {
		return 0, 0, 0, false, 0, false, false, fmt.Errorf("seed must be non-negative")
	}

	// Parse withClues parameter
	withCluesStr := c.Query("withClues", "false")
	withClues, err = strconv.ParseBool(withCluesStr)
	if err != nil {
		return 0, 0, 0, false, 0, false, false, fmt.Errorf("invalid withClues parameter: %w", err)
	}

	// Parse cluesSeed parameter
	cluesSeedStr := c.Query("cluesSeed", "0") // Default to 0 (random seed) if not provided
	cluesSeed, err = strconv.ParseInt(cluesSeedStr, 10, 64)
	if err != nil {
		return 0, 0, 0, false, 0, false, false, fmt.Errorf("invalid cluesSeed parameter: %w", err)
	}
	if !isSeedValid(cluesSeed) {
		return 0, 0, 0, false, 0, false, false, fmt.Errorf("cluesSeed must be non-negative")
	}

	// Parse unsolved parameter
	unsolvedStr := c.Query("unsolved", "false")
	unsolved, err = strconv.ParseBool(unsolvedStr)
	if err != nil {
		return 0, 0, 0, false, 0, false, false, fmt.Errorf("invalid unsolved parameter: %w", err)
	}

	// Parse cryptic parameter (not used in this version, but kept for future use)
	crypticStr := c.Query("cryptic", "false")
	cryptic, err = strconv.ParseBool(crypticStr)
	if err != nil {
		return 0, 0, 0, false, 0, false, false, fmt.Errorf("invalid cryptic parameter: %w", err)
	}

	return rows, cols, seed, withClues, cluesSeed, unsolved, cryptic, nil
}

// prepareCrosswordResponse prepares the JSON response for the crossword
func prepareCrosswordResponse(crosswordResult crossword.CrosswordResult, cluesResult *cluer.CluesResult, unsolved bool) CrosswordResponse {
	cw := crosswordResult.Crossword

	// Extract the data from the crossword
	data := extractData(cw, unsolved)

	// Extract the words
	words := extractWords(cw, cluesResult, unsolved)

	// Prepare the response
	response := CrosswordResponse{
		CrosswordSeed: crosswordResult.Seed,
		Rows:          cw.Rows(),
		Cols:          cw.Columns(),
		Data:          data,
		Words:         words,
	}

	// Add clues seed if clues were generated
	if cluesResult != nil {
		response.CluesSeed = cluesResult.Seed
	}

	return response
}

func extractData(cw *crossword.Crossword, unsolved bool) string {
	data := make([]byte, 0, cw.Columns()*cw.Rows())

	for letter := crossword.CrosswordLetter(cw); letter != nil; letter = letter.Next() {
		if !letter.IsBlank() && unsolved {
			data = append(data, byte('*'))
		} else {
			data = append(data, letter.GetValue())
		}
	}
	return string(data)
}

// extractWords extracts all words from the crossword
func extractWords(cw *crossword.Crossword, cluesResult *cluer.CluesResult, unsolved bool) []CrosswordWord {
	var words []CrosswordWord

	// Start with horizontal words
	for word := crossword.Word(cw); word != nil; word = word.Next() {
		direction := word.GetDirection().String()
		pos := word.GetPos()
		value := string(word.GetValue())

		wordRef := CrosswordWord{
			Pos:       pos,
			Direction: direction,
		}

		// Only set value if not unsolved
		if !unsolved {
			wordRef.Value = value
		}

		// Add clue if available
		if cluesResult != nil {
			if clue, exists := cluesResult.Clues[value]; exists {
				wordRef.Clue = clue
			}
		}

		words = append(words, wordRef)
	}

	return words
}

// isSizeValid checks if the size is valid (between 3 and 15)
func isSizeValid(size int) bool {
	return size >= 3 && size <= 15
}

// isSeedValid checks if the seed is valid (non-negative)
func isSeedValid(seed int64) bool {
	return seed >= 0
}
