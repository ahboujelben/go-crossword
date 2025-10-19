# 🎮 Go-Crossword MCP Server

![Go](https://img.shields.io/badge/Go-1.24%2B-blue)
![MCP](https://img.shields.io/badge/MCP-1.0.0-green)
![Docker](https://img.shields.io/badge/Docker-Supported-blue)
![License](https://img.shields.io/badge/License-MIT-green)

> 🧩 **Generate engaging crossword puzzles directly from your AI assistant!**

This [MCP](https://modelcontextprotocol.io) server brings powerful crossword puzzle generation capabilities to AI assistants like Claude, C Zed, and other MCP-compatible clients.

---

## ✨ Features

- 🎲 **Dynamic Puzzle Generation** - Create crossword puzzles from 3x3 up to 15x15 grids
- 🤖 **AI-Native Integration** - Seamlessly integrates with MCP-compatible AI assistants
- 🎯 **Smart Word Placement** - Advanced algorithm ensures proper word intersections
- 📚 **Curated Dictionary** - Generates genuinely interesting crosswords!
- 📝 **Intelligent Clue Generation** - Delegates clue creation to your AI assistant for contextual, creative hints
- 🚀 **Ultra-Lightweight** - Built with Go and packaged in a minimal scratch container (~5MB)
- ⚡ **Fast & Efficient** - Multi-threaded generation with optimized performance
- 🔒 **Secure** - No external dependencies in runtime, runs in isolated container

---

## 🚀 Quick Start

### Pull the Image

```bash
docker pull ahboujelben/go-crossword-mcp
```

### Test the Server

```bash
docker run --rm -i ahboujelben/go-crossword-mcp
```

---

## 🔌 Integration with AI Assistants

### Claude Desktop

Add to your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "go-crossword": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "ahboujelben/go-crossword-mcp"]
    }
  }
}
```

### VSCode Editor

Add to your VSCode MCP servers (`mcp.json`):

```json
{
  "servers": {
    "go-crossword": {
      "type": "stdio",
      "command": "docker",
      "args": ["run", "--rm", "-i", "ahboujelben/go-crossword-mcp"]
    }
  },
}
```

### Zed Editor

Add to your Zed settings (`settings.json`):

```json
{
  "context_servers": {
    "go-crossword": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "ahboujelben/go-crossword-mcp"]
    }
  }
}
```

### Other MCP Clients

The server communicates via stdio using the MCP protocol. Configure your client to run:

```bash
docker run --rm -i ahboujelben/go-crossword-mcp
```

---

## 🎯 Usage

Once integrated, your AI assistant gains access to the `generate-crossword` tool. Simply ask:

- *"Generate a 10x10 crossword puzzle"*
- *"Create a small crossword for me"*
- *"Make a 7x5 crossword with cryptic clues"*

The AI will:
1. 📊 Generate the crossword grid structure
2. 🎨 Display the unsolved puzzle in a clean format
3. 💡 Create engaging clues for each word (Across & Down)
4. ✅ Keep the solution hidden until you request it

---

## 🛠️ Tool Specification

### `generate-crossword`

Generates a crossword puzzle with specified dimensions.

**Input Parameters:**
- `rows` (int): Number of rows (3-15)
- `cols` (int): Number of columns (3-15)

**Output:**
- `unsolvedCrossword` (string): The puzzle grid without solutions
- `solvedCrossword` (string): The complete puzzle with all answers
- `rowWords` (array): List of horizontal words with positions
- `columnWords` (array): List of vertical words with positions

---

## 🏗️ Architecture

Built with modern Go and leveraging:
- **Model Context Protocol SDK** - Official Go implementation
- **Sophisticated Word Placement Engine** - Optimized for quality intersections
- **Multi-threaded Generation** - Fast puzzle creation
- **Compact Rendering** - Beautiful terminal-friendly output

---

## 📦 Image Details

- **Base Image**: `scratch` (minimal, security-focused)
- **Size**: ~5MB
- **Architecture**: linux/amd64
- **Go Version**: 1.24+
- **Entry Point**: Stdio-based MCP server

---

## 🔗 Links

- 📚 **GitHub Repository**: [github.com/ahboujelben/go-crossword](https://github.com/ahboujelben/go-crossword)
- 🌐 **MCP Protocol**: [modelcontextprotocol.io](https://modelcontextprotocol.io)
- 📖 **Documentation**: [Full README](https://github.com/ahboujelben/go-crossword#readme)

---

## 📄 License

MIT License - Free for personal and commercial use

---

## 🙌 Contributing

Contributions, issues, and feature requests are welcome!
Visit the [GitHub repository](https://github.com/ahboujelben/go-crossword) to get involved.

---

**Made with ❤️ by ahboujelben** | Powered by Go & MCP
