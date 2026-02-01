package config

import "fmt"

var (
	Config RootConfig
)

type ChatModelConfig struct {
	Provider string `toml:"provider"`
	BaseURL  string `toml:"base_url"`
	Model    string `toml:"model"`
}

type EmbeddingConfig struct {
	Provider string `toml:"provider"`
	BaseURL  string `toml:"base_url"`
	Model    string `toml:"model"`
}

type VectorStoreConfig struct {
	Provider   string `toml:"provider"`
	BaseURL    string `toml:"base_url"`
	Collection string `toml:"collection"`
}

type RootConfig struct {
	ChatModelConfig   *ChatModelConfig   `toml:"chat_model"`
	EmbeddingConfig   *EmbeddingConfig   `toml:"embeddings"`
	VectorstoreConfig *VectorStoreConfig `toml:"vector_store"`
}

func (c *RootConfig) Validate() error {
	validModelProviders := map[string]bool{
		"ollama":   true,
		"llamacpp": true,
		"openai":   true,
	}

	validVectorStoreProviders := map[string]bool{
		"qdrant": true,
	}

	// [chat_model] validation
	if c.ChatModelConfig == nil {
		return fmt.Errorf("Missing chat_model config.")
	}
	if c.ChatModelConfig.Model == "" {
		return fmt.Errorf("Invalid chat_model.model: model name is required.")
	}
	if c.ChatModelConfig.BaseURL == "" {
		return fmt.Errorf("Invalid chat_model.base_url: base_url is required.")
	}
	if !validModelProviders[c.ChatModelConfig.Provider] {
		return fmt.Errorf("Invalid chat model provider: \"%s\". Valid providers include: \"ollama\", \"llama.cpp\", \"openai\"", c.ChatModelConfig.Provider)
	}

	// [embeddings] validation
	if c.EmbeddingConfig == nil {
		return fmt.Errorf("Missing embeddings config.")
	}
	if c.EmbeddingConfig.Model == "" {
		return fmt.Errorf("Invalid embeddings.model: model name is required.")
	}
	if c.EmbeddingConfig.BaseURL == "" {
		return fmt.Errorf("Invalid embeddings.base_url: base_url is required.")
	}
	if !validModelProviders[c.EmbeddingConfig.Provider] {
		return fmt.Errorf("Invalid embedding provider: \"%s\". Valid providers include: \"ollama\", \"llama.cpp\", \"openai\"", c.EmbeddingConfig.Provider)
	}

	// [vector_store] validation
	if c.VectorstoreConfig == nil {
		return fmt.Errorf("Missing vector_store config.")
	}
	if c.VectorstoreConfig.Collection == "" {
		return fmt.Errorf("Invalid vector_store.collection: collection name is required.")
	}
	if c.VectorstoreConfig.BaseURL == "" {
		return fmt.Errorf("Invalid vector_store.base_url: base_url is required.")
	}
	if !validVectorStoreProviders[c.VectorstoreConfig.Provider] {
		return fmt.Errorf("Invalid vector store provider: \"%s\". Valid providers include: \"qdrant\"", c.VectorstoreConfig.Provider)
	}

	return nil
}
