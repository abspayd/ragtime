package config

import "fmt"

var (
	Config RootConfig
)

type ChatModelConfig struct {
	Provider string `toml:"provider"`
	Model    string `toml:"model"`
}

type EmbeddingConfig struct {
	Provider string `toml:"provider"`
	Model    string `toml:"model"`
}

type VectorStoreConfig struct {
	Provider   string `toml:"provider"`
	Collection string `toml:"collection"`
}

type RootConfig struct {
	Chat_Model_Config  *ChatModelConfig   `toml:"chat_model"`
	Embedding_Config   *EmbeddingConfig   `toml:"embeddings"`
	Vectorstore_Config *VectorStoreConfig `toml:"vector_store"`
}

func (c *RootConfig) Validate() error {
	validModelProviders := map[string]bool{
		"ollama":   true,
		"llamacpp": true,
	}

	validVectorStoreProviders := map[string]bool{
		"qdrant": true,
	}

	// [chat_model] validation
	if c.Chat_Model_Config == nil {
		return fmt.Errorf("Missing chat_model config.")
	}

	if c.Chat_Model_Config.Model == "" {
		return fmt.Errorf("Invalid chat_model.model: model name is required.")
	}

	if !validModelProviders[c.Chat_Model_Config.Provider] {
		return fmt.Errorf("Invalid chat model provider: \"%s\". Valid providers include: \"ollama\", \"llama.cpp\"", c.Chat_Model_Config.Provider)
	}

	// [embeddings] validation
	if c.Embedding_Config == nil {
		return fmt.Errorf("Missing embeddings config.")
	}

	if c.Embedding_Config.Model == "" {
		return fmt.Errorf("Invalid embeddings.model: model name is required.")
	}

	if !validModelProviders[c.Embedding_Config.Provider] {
		return fmt.Errorf("Invalid embedding provider: \"%s\". Valid providers include: \"ollama\", \"llama.cpp\"", c.Embedding_Config.Provider)
	}

	// [vector_store] validation
	if c.Vectorstore_Config == nil {
		return fmt.Errorf("Missing vector_store config.")
	}

	if c.Vectorstore_Config.Collection == "" {
		return fmt.Errorf("Invalid vector_store.collection: collection name is required.")
	}

	if !validVectorStoreProviders[c.Vectorstore_Config.Provider] {
		return fmt.Errorf("Invalid vector store provider: \"%s\". Valid providers include: \"qdrant\"", c.Vectorstore_Config.Provider)
	}

	return nil
}
