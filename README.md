# Ragtime

Local RAG toolkit for document Q&A. Ingest documents, query with natural language, and chat with your local LLMs.

## Status

This project is under active development. See [PLAN.md](PLAN.md) for the roadmap.

## Features

- [x] CLI framework and config loading
- [ ] Document ingestion (WIP)
- [ ] Semantic search
- [x] LLM integration (llamacpp, ollama)
- [ ] Vector storage (Qdrant)
- [ ] Interactive chat TUI

## Installation

```bash
git clone https://github.com/abspayd/ragtime.git
cd ragtime
go build -o ragtime ./cmd/ragtime
```

### Dependencies

- [Qdrant](https://qdrant.tech/) - Vector database
- [llamacpp](https://github.com/ggerganov/llama.cpp) or [Ollama](https://ollama.ai/) - Local LLM inference
- `pdftotext` (optional) - For PDF extraction (`poppler-utils` package)

## Quick Start

1. Copy and edit the config:
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

2. Build and run:
   ```bash
   go build -o ragtime ./cmd/ragtime
   ./ragtime --help
   ```

Additional steps (once implemented):

3. Start Qdrant:
   ```bash
   docker compose up -d
   ```

4. Ingest documents:
   ```bash
   ragtime ingest ./docs --recursive
   ```

5. Query or chat:
   ```bash
   ragtime query "What is the main topic?"
   ragtime chat
   ```

## Configuration

### config.toml

```toml
[vector_store]
provider = "qdrant"
collection = "documents"

[chat_model]
provider = "llamacpp"  # or "ollama"
model = "gpt-oss"

[embeddings]
provider = "llamacpp"  # or "ollama"
model = "nomic-embed-text-v1.5"
```

### .env

```bash
qdrant_url="http://localhost:6333"
qdrant_api_key="MY_KEY"

ollama_url="http://localhost:11434"
llamacpp_url="http://localhost:8081/v1"
```

## Commands

| Command | Description | Status |
|---------|-------------|--------|
| `ragtime --help` | Show available commands | Done |
| `ragtime ingest <path>` | Add documents to the vector store | Planned |
| `ragtime query "<question>"` | Ask a question | Planned |
| `ragtime chat` | Interactive chat REPL | Planned |
| `ragtime serve` | Start HTTP API server | Planned |
| `ragtime collections` | Manage collections | Planned |
| `ragtime docs` | Manage indexed documents | Planned |

## License

MIT
