# ragtime

A RAG utility for local document and code queries

`ragtime` ingests your files locally into a Qdrant database, generates embeddings with an OpenAI-compatible model server, and answers questions using the stored document context.

## Features

- Ingest one or more local files into a Qdrant collection
- Generate embeddings with llama.cpp, Ollama, or OpenAI-compatible endpoints
- Improve the answering capabilities of your local models using your local datasets
- Uses tree-sitter to chunk code and markdown source files into pieces queryable by your local LLMs
- TOML config and `.env` support

## Install

**Prerequisites:** Go 1.25.6 or later ([golang.org/dl](https://go.dev/dl/)), [Docker](https://www.docker.com/) for Qdrant, and a running OpenAI-compatible chat and embeddings server

```sh
git clone https://github.com/abspayd/ragtime.git
cd ragtime
go build -o ragtime ./cmd/ragtime
```

## Usage

Start Qdrant:

```sh
make up
```
Ingest files into the Qdrant database:

```sh
./ragtime ingest path/to/file.md path/to/source.go
```

Ask a question:

```sh
./ragtime chat "What does this project do?"
```

| Command | Description |
|---------|-------------|
| `ingest path...` | Add files to the vector store |
| `chat message...` | Chat with an LLM using RAG search |

**Global flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `-c, --config` | `config.toml` | Path to the config file |
| `-l, --log` | `logs/ragtime.log` | Path to the log file |
| `-v, --verbose` | `false` | Include more details in the logs |

**`ingest` flags:**

| Flag | Default | Description |
|------|---------|-------------|
| `--no-git-ignore` | `false` | Do not ignore files in `.gitignore` |
| `-I, --ignore` | | Glob pattern to exclude files or directories |

## Config

```toml
[vector_store]
provider = "qdrant"
base_url = "localhost"
collection = "documents"

[chat_model]
provider = "llamacpp"
base_url = "http://localhost:8080"
model = "ministral-3-14b"

[embeddings]
provider = "llamacpp"
base_url = "http://localhost:8080"
model = "nomic-embed-text-v1.5"
```

Optional environment variables:

```sh
QDRANT_API_KEY="MY_KEY"
OPENAI_API_KEY="MY_KEY"
```

## Notes

- Qdrant runs on ports `6333` and `6334` when started with `make up`.
- The ingest command recreates stored chunks for a file path before uploading new chunks.
- Directory ingestion and PDF extraction are not implemented yet.

## Future plans

- Add an interactive chat TUI
- Add agent and tool support
