# Ragtime

Local RAG toolkit for document Q&A. Ingest documents, query with natural language, and chat with your local LLMs.

## Features

- **Document ingestion** - Index PDFs, text files, and markdown
- **Semantic search** - Query your documents with natural language
- **Local LLMs** - Works with llamacpp, ollama, or OpenAI-compatible APIs
- **Vector storage** - Uses Qdrant for fast similarity search
- **Interactive chat** - TUI interface for conversational Q&A

## Installation

```bash
go install github.com/abspayd/ragtime/cmd/ragtime@latest
```

Or build from source:

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

1. Start Qdrant:
   ```bash
   docker compose up -d
   ```

2. Copy and edit the config:
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

3. Ingest documents:
   ```bash
   ragtime ingest ./docs --recursive
   ```

4. Query:
   ```bash
   ragtime query "What is the main topic?"
   ```

5. Or start an interactive chat:
   ```bash
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

| Command | Description |
|---------|-------------|
| `ragtime ingest <path>` | Add documents to the vector store |
| `ragtime query "<question>"` | Ask a question |
| `ragtime chat` | Interactive chat REPL |
| `ragtime serve` | Start HTTP API server |
| `ragtime collections` | Manage collections |
| `ragtime docs` | Manage indexed documents |

## License

MIT
