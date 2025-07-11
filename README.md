# 🎮 GoCrossword

![Status](https://img.shields.io/badge/Status-Active-success)
![Go](https://img.shields.io/badge/Go-1.24%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)

## 🧩 Overview

GoCrossword is a powerful crossword toolkit that creates engaging crossword puzzles from scratch! The system fills an empty grid with words from a predefined dictionary and can generate clever clues using AI.

**✨ Features:**

- 🎲 Create random or seeded crossword puzzles
- 🤖 AI-powered clue generation with [Ollama](https://github.com/ollama/ollama)
- 🖥️ CLI tool for quick puzzle creation
- 🌐 REST API for web/mobile integration
- 🐳 Docker support for easy deployment

## 🚀 Quick Start

### Prerequisites

- Go 1.24 or higher
- [Ollama](https://github.com/ollama/ollama) running locally with `llama3:8b` model (for clue generation)

### Installation

```shell
# Clone the repository
git clone https://github.com/ahboujelben/go-crossword.git
cd go-crossword

# Build the CLI tool
make build-cli
```

## 💻 CLI Usage

The CLI tool allows you to quickly generate crossword puzzles for printing or sharing.

```shell
# Run using the CLI executable
./go-crossword-cli

# Or use the Make target
make run-cli

# Docker option
make docker-run-cli
```

### CLI Options

```shell
Usage: go-crossword-cli [options]

Options:
  -rows int            Number of rows in the crossword grid (default 13)
  -cols int            Number of columns in the crossword grid (default 13)
  -crossword-seed int  Seed for crossword generation (default: random)
  -clues-seed int      Seed for clue generation (default: random)
  -unsolved            Display puzzle in unsolved mode
  -cryptic             Generate more cryptic/challenging clues
  -ollama-url string   URL of the Ollama server (default "http://localhost:11434")
  -ollama-model string Model to use for generating clues (default "llama3:8b")
  -compact             Use a more compact rendering style
```

### CLI Examples

#### Generate a random 13x13 crossword grid

![Random crossword grid example](vhs/1-plain-crossword.gif)

Note the seed printed at the end can be used to generate clues for that specific crossword.

#### Generate a random crossword grid with custom dimensions

![Custom dimension crossword example](vhs/2-plain-crossword-custom-size.gif)

#### Generate a crossword with clues in solved mode

![Crossword with clues example](vhs/3-crossword-with-clues.gif)

The example above shows how clues are generated for a specific seeded crossword. If the seed is omitted, a random crossword grid is created along with the clues.

Note that the second seed `-clues-seed` allows recreating that crossword later with those exact clues.

#### Cryptic clues

By default, the clues should yield a crossword of normal/easy difficulty. Passing `-cryptic` will prompt the LLM model to come up with more cryptic clues.

![Cryptic clues example](vhs/4-crossword-with-clues-cryptic.gif)

#### Unsolved mode

Passing `-unsolved` will hide the solution of the generated crossword. This can apply to previously generated crosswords (when the crossword/clues seed values are passed) or a new random one.

![Unsolved crossword example](vhs/5-crossword-with-clues-unsolved.gif)

If a random unsolved crossword is generated, the solution can be shown by rerunning the command with the printed seeds.

![Solved crossword example](vhs/6-crossword-with-clues-unsolved-2.gif)

## 🌐 API Server

GoCrossword also provides a RESTful API server that allows you to generate crosswords programmatically, perfect for web applications or mobile apps!

```shell
# Start the API server
make run-api

# Or use Docker
make docker-run-api

# Start both API and Ollama using Docker Compose
make docker-compose-up

# Example API request to generate a crossword
curl -s "http://localhost:8080/api/crossword?rows=7&cols=7"
```

### API Endpoints

#### Generate a Crossword

```http
GET /api/crossword
```

Query parameters:

- `rows`: Number of rows (default: 13)
- `cols`: Number of columns (default: 13)
- `crosswordSeed`: Seed for crossword generation (optional)
- `cluesSeed`: Seed for clue generation (optional)
- `withClues`: Include clues in the response (default: false)
- `cryptic`: Generate more cryptic clues (default: false)
- `unsolved`: Return the crossword in unsolved mode (default: false)

#### Example Response

```json
{
  "crosswordSeed": 4453768179316043008,
  "cluesSeed": 4706815147906524866,
  "rows": 7,
  "cols": 7,
  "data": ".maniacw.m.n.has.asiar.m.h.smeadowse.l.r.idrivers",
  "words": [
    {
      "pos": 1,
      "direction": "horizontal",
      "value": "maniac",
      "clue": "Frenzied devotee of intense activity"
    },
    {
      "pos": 14,
      "direction": "horizontal",
      "value": "as",
      "clue": "Used in equal quantities"
    },
    {
      "pos": 17,
      "direction": "horizontal",
      "value": "asia",
      "clue": "Continent often referred to as the Far East"
    },
    {
      "pos": 28,
      "direction": "horizontal",
      "value": "meadows",
      "clue": "Grassy areas in rolling countryside"
    },
    {
      "pos": 42,
      "direction": "horizontal",
      "value": "drivers",
      "clue": "Automotive professionals are typically licensed"
    },
    {
      "pos": 7,
      "direction": "vertical",
      "value": "warmed",
      "clue": "Heated to a higher temperature"
    },
    {
      "pos": 2,
      "direction": "vertical",
      "value": "am",
      "clue": "Formal mode of address"
    },
    {
      "pos": 23,
      "direction": "vertical",
      "value": "mali",
      "clue": "Landlocked West African country where the Niger River flows"
    },
    {
      "pos": 4,
      "direction": "vertical",
      "value": "inshore",
      "clue": "Fisherman's zone where waves are calm"
    },
    {
      "pos": 6,
      "direction": "vertical",
      "value": "chassis",
      "clue": "Underlying structure for cars and planes"
    }
  ]
}
```

### Configuration

The API server can be configured using environment variables:

- `PORT`: Server port (default: 8080)
- `OLLAMA_URL`: URL of the Ollama server (default: `http://localhost:11434`)
- `OLLAMA_MODEL`: Model to use for clue generation (default: llama3:8b)

## 🤖 Ollama Integration

GoCrossword uses [Ollama](https://github.com/ollama/ollama) to generate clever and engaging clues for your crosswords.

### Setup Ollama

1. Install Ollama from [https://ollama.com/](https://ollama.com/)
2. Pull the default model:

   ```shell
   ollama pull llama3:8b
   ```

3. Start the Ollama server:

   ```shell
   ollama serve
   ```

### Docker Setup

The included Docker Compose file automatically sets up Ollama with the required model:

```yaml
services:
  api:
    # Configuration for the API service
    depends_on:
      - ollama
    environment:
      - OLLAMA_URL=http://ollama:11434
      
  ollama:
    image: ollama/ollama:latest
    volumes:
      - ollama_data:/root/.ollama
    entrypoint: >
      sh -c "ollama serve &
             sleep 10 &&
             ollama pull llama3:8b &&
             wait"
```

## 🧠 How It Works

GoCrossword uses a sophisticated algorithm to generate crossword puzzles:

1. **Grid Generation**: Creates a grid of the specified dimensions
2. **Word Placement**: Places words from a dictionary into the grid, ensuring proper intersections
3. **Clue Generation**: Uses Ollama's LLM capabilities to create engaging clues for each word
4. **Rendering**: Outputs the crossword in various formats (CLI text, JSON for API)

### Architecture Diagram

```mermaid
graph TB
    subgraph "CLI Application"
        CLI[CLI Entry Point]
        CR[Crossword Renderer]
    end

    subgraph "API Server"
        API[API Server]
        REST[REST Endpoints]
    end

    subgraph "Core Modules"
        CW[Crossword Generator]
        DICT[Dictionary]
        CLUE[Clue Generator]
    end

    subgraph "External Dependencies"
        OLLAMA[Ollama LLM]
    end

    CLI --> CW
    CLI --> CR
    API --> REST
    REST --> CW
    CW --> DICT
    CW --> CLUE
    CLUE --> OLLAMA

    classDef core fill:#f9f,stroke:#333,stroke-width:2px;
    classDef ext fill:#bbf,stroke:#333,stroke-width:2px;
    classDef app fill:#bfb,stroke:#333,stroke-width:2px;
    
    class CW,DICT,CLUE core;
    class OLLAMA ext;
    class CLI,CR,API,REST app;
```

## 🛠️ Development

### Project Structure

```text
go-crossword/
├── api/           # API server implementation
├── cli/           # Command-line interface
├── modules/       # Core modules (crossword, dictionary, clue generation)
├── docker-compose.yml
└── Makefile       # Build and run targets
```

### Available Make Commands

```bash
make build-cli         # Build the CLI application
make run-cli           # Run the CLI application
make build-api         # Build the API server
make run-api           # Run the API server
make test              # Run tests
make docker-build-cli  # Build CLI Docker image
make docker-run-cli    # Run CLI Docker container
make docker-build-api  # Build API Docker image
make docker-run-api    # Run API Docker container
make docker-compose-up # Start all services with Docker Compose
```

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙌 Contributing

Contributions are welcome! Feel free to open issues or submit pull requests.
