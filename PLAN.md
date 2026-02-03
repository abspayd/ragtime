# Project plan

Build a RAG application in Go with a CLI interface, HTTP API, and TUI chat.

## Project Structure

```
ragtime/
├── cmd/ragtime/main.go           # Entry point
├── internal/
│   ├── config/                   # TOML + env config
│   ├── cli/                      # Cobra commands
│   ├── llm/                      # OpenAI-compatible client
│   ├── embedding/                # Embedding client
│   ├── vectorstore/qdrant/       # Qdrant client
│   ├── document/                 # Chunking + ingest pipeline
│   ├── rag/                      # Retriever, engine
│   ├── agent/                    # State machine
│   ├── server/                   # HTTP API
│   └── tui/                      # Bubble Tea chat
├── config.toml                   # App config (keep existing)
└── .env                          # Secrets (keep existing)
```

---

## Phase 1: Foundation (Start Here)

**Goal**: CLI skeleton + config system working

### Files to Create

| File | Purpose |
|------|---------|
| `cmd/ragtime/main.go` | Entry point, calls `cli.Execute()` |
| `internal/config/types.go` | Config structs matching existing TOML |
| `internal/config/config.go` | Load TOML + .env, validate |
| `internal/cli/root.go` | Cobra root + global flags |
| `internal/cli/init.go` | `ragtime init` - create default config |

### Config Types

```go
type Config struct {
    VectorStore VectorStoreConfig `toml:"vector_store"`
    ChatModel   ChatModelConfig   `toml:"chat_model"`
    Embeddings  EmbeddingConfig   `toml:"embeddings"`
}

type Env struct {
    QdrantURL    string // qdrant_url
    QdrantAPIKey string // qdrant_api_key
    OllamaURL    string // ollama_url
    LlamaCPPURL  string // llamacpp_url
}
```

### Dependencies

```
github.com/spf13/cobra
github.com/BurntSushi/toml
github.com/joho/godotenv
```

### Verification

- [x] `go build ./cmd/ragtime` succeeds
- [x] `ragtime --help` shows available commands
- [x] `ragtime init` creates config.toml template
- [x] Config loads from existing `config.toml`
- [x] Env loads from existing `.env`

---

## Phase 2: Core Services

**Goal**: LLM, embedding, and Qdrant clients working

### Files to Create

| File | Purpose |
|------|---------|
| `internal/llm/types.go` | Message, ChatRequest/Response types |
| `internal/llm/client.go` | Client interface |
| `internal/llm/openai.go` | OpenAI-compatible HTTP client |
| `internal/embedding/client.go` | Embedding interface + impl |
| `internal/vectorstore/types.go` | Document, SearchResult types |
| `internal/vectorstore/qdrant/client.go` | Qdrant client |


### Dependencies

```
github.com/qdrant/go-client
```

### Verification

- [x] LLM client gets response from llamacpp
- [ ] LLM streaming works (SSE parsing)
- [x] Embedding client generates vectors
- [ ] Qdrant: create collection, upsert, search

---

## Phase 3: Document Pipeline

**Goal**: `ragtime ingest` working end-to-end

### Files to Create

| File | Purpose |
|------|---------|
| `internal/document/types.go` | Chunk type + metadata |
| `internal/document/pipeline.go` | Orchestrates read→chunk→embed→store |
| `internal/document/text_chunker.go` | Text splitter (prose) |
| `internal/document/code_chunker.go` | Tree-sitter code chunker |
| `internal/document/queries/*.scm` | Per-language tree-sitter queries |
| `internal/cli/ingest.go` | `ragtime ingest` command |

### File Loading

File loading with PDF support.

```go
func (p *Pipeline) readFile(ctx context.Context, path string) ([]byte, error) {
    if filepath.Ext(path) == ".pdf" {
        return p.extractPDF(ctx, path)
    }
    return os.ReadFile(path)
}

func (p *Pipeline) extractPDF(ctx context.Context, path string) ([]byte, error) {
    return exec.CommandContext(ctx, "pdftotext", "-layout", path, "-").Output()
}
```

### Code Chunking (tree-sitter with fallback)

Use tree-sitter queries to extract semantic chunks from code. Each language gets a `.scm` query file—no new Go code needed per language.

**Query files define chunk boundaries:**

```scheme
;; queries/go.scm
(function_declaration) @chunk
(method_declaration) @chunk
(type_declaration) @chunk

;; queries/python.scm
(function_definition) @chunk
(class_definition) @chunk
```

**Generic chunker with fallback:**

```go
type CodeChunker struct {
    queries map[string][]byte  // language -> query file contents
}

func (c *CodeChunker) Chunk(lang string, content []byte) ([]Chunk, error) {
    query, ok := c.queries[lang]
    if !ok {
        // Fallback: text splitting for unknown languages
        return c.fallbackChunk(content)
    }
    return c.treesitterChunk(lang, content, query)
}
```

**Benefits:**
- Add new language support by dropping in a `.scm` file
- Chunks align with semantic boundaries (functions, classes)
- Falls back gracefully for unsupported languages

### Verification

- [ ] `ragtime ingest ./file.txt` adds to Qdrant
- [ ] `ragtime ingest ./doc.pdf` extracts and indexes
- [ ] `ragtime ingest ./docs/ --recursive` processes directory
- [ ] Chunks have correct metadata (source, index)

---

## Phase 4: RAG Engine

**Goal**: `ragtime query` working

### Files to Create

| File | Purpose |
|------|---------|
| `internal/rag/retriever.go` | Embed query → search Qdrant |
| `internal/rag/prompts.go` | System prompt, context template |
| `internal/rag/engine.go` | Retrieve → build prompt → generate |
| `internal/cli/query.go` | `ragtime query` command |

### Verification

- [ ] `ragtime query "question"` returns answer with sources
- [ ] `ragtime query --no-llm "keyword"` returns raw chunks
- [ ] `ragtime query -v` shows source citations
- [ ] Streaming works (`ragtime query --stream`)

---

## Phase 5: Agent System

**Goal**: Multi-turn agentic conversations with tool use

### Files to Create

| File | Purpose |
|------|---------|
| `internal/agent/state.go` | State enum, Context type |
| `internal/agent/tool.go` | Tool interface, ToolCall types |
| `internal/agent/agent.go` | State machine loop |
| `internal/agent/tools/search.go` | Built-in search tool |

### State Machine

```
idle → thinking → [tool_call → thinking]* → complete
                ↘ error
```

### Verification

- [ ] Agent answers without tools when context sufficient
- [ ] Agent invokes search tool when needed
- [ ] Tool results incorporated into response
- [ ] Max turns limit prevents infinite loops

---

## Phase 6: Interface Layer

**Goal**: HTTP API + TUI chat

### Files to Create

| File | Purpose |
|------|---------|
| `internal/server/server.go` | HTTP routes, handlers |
| `internal/server/handlers.go` | Query, ingest, collections endpoints |
| `internal/tui/model.go` | Bubble Tea model |
| `internal/tui/styles.go` | Lipgloss styles |
| `internal/tui/chat.go` | Chat view + streaming |
| `internal/cli/serve.go` | `ragtime serve` command |
| `internal/cli/chat.go` | `ragtime chat` command |

### HTTP Endpoints

```
POST /api/query         - Query with RAG
POST /api/query/stream  - SSE streaming
POST /api/ingest        - Add document
GET  /api/collections   - List collections
GET  /health            - Health check
```

### TUI Dependencies

```
github.com/charmbracelet/bubbletea
github.com/charmbracelet/bubbles
github.com/charmbracelet/lipgloss
github.com/charmbracelet/glamour  # Markdown rendering
```

### Verification

- [ ] `ragtime serve` starts on :8080
- [ ] `curl /api/query` returns JSON response
- [ ] SSE streaming works in browser
- [ ] `ragtime chat` launches TUI
- [ ] TUI streams responses in real-time
- [ ] Ctrl+C exits cleanly

---

## Implementation Order

| Phase | Focus | Deliverable |
|-------|-------|-------------|
| **1** | Foundation | `ragtime --help`, config loading |
| **2** | Core Services | LLM/embedding/Qdrant working |
| **3** | Document Pipeline | `ragtime ingest` working |
| **4** | RAG Engine | `ragtime query` working |
| **5** | Agent System | Multi-turn conversations |
| **6** | Interface Layer | HTTP API + TUI chat |

**Recommended first milestone**: Complete phases 1-4 for a minimal working RAG system, then iterate on 5-6.

---

## External Dependencies Summary

```go
require (
    // CLI
    github.com/spf13/cobra v1.8.0

    // Config
    github.com/BurntSushi/toml
    github.com/joho/godotenv v1.5.1

    // Vector Store
    github.com/qdrant/go-client v1.10.0

    // TUI (Phase 6)
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/bubbles v0.18.0
    github.com/charmbracelet/lipgloss v0.9.0
    github.com/charmbracelet/glamour v0.6.0
)
```

---

## Notes

- Keep existing `config.toml` and `.env.example` formats
- PDF extraction uses system `pdftotext` (poppler-utils) - no CGO needed
- TypeScript scraper can be added later as separate service in `extractors/scraper/`
