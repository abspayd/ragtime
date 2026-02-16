package cli

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/abspayd/ragtime/internal/cli/chat"
	"github.com/abspayd/ragtime/internal/cli/ingest"
	"github.com/abspayd/ragtime/internal/config"
	"github.com/abspayd/ragtime/internal/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
	logFile string

	rootCmd = &cobra.Command{
		Use:   "ragtime",
		Short: "Local RAG toolkit for document Q&A",
		Long: `
██████╗  █████╗  ██████╗████████╗██╗███╗   ███╗███████╗
██╔══██╗██╔══██╗██╔════╝╚══██╔══╝██║████╗ ████║██╔════╝
██████╔╝███████║██║  ███╗  ██║   ██║██╔████╔██║█████╗  
██╔══██╗██╔══██║██║   ██║  ██║   ██║██║╚██╔╝██║██╔══╝  
██║  ██║██║  ██║╚██████╔╝  ██║   ██║██║ ╚═╝ ██║███████╗
╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝   ╚═╝   ╚═╝╚═╝     ╚═╝╚══════╝

Ingest documents, query with natural language, and chat with your local LLMs.
		`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := initLogs(); err != nil {
				return err
			}
			if err := initConfig(); err != nil {
				return err
			}

			return nil
		},
		// RunE: func(cmd *cobra.Command, args []string) error {

		// 	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		// 	defer cancel()

		// 	chatClient := models.NewOpenAIClient("gpt-oss", config.Config.ChatModelConfig.BaseURL, "")
		// 	messages, err := chatClient.Chat(ctx, []models.Message{
		// 		{
		// 			Role:    "user",
		// 			Content: "Hello",
		// 		},
		// 	})
		// 	if err != nil {
		// 		return err
		// 	}

		// 	for _, message := range messages {
		// 		fmt.Printf("%s: %s\n", message.Role, message.Content)
		// 	}

		// 	embeddingsClient := models.NewOpenAIClient(config.Config.EmbeddingConfig.Model, config.Config.EmbeddingConfig.BaseURL, "")

		// 	embeddingsResp, err := embeddingsClient.Embed(ctx, "Hello, this is a test for how embeddings work")
		// 	if err != nil {
		// 		return err
		// 	}

		// 	qdrant.NewClient(&qdrant.Config{
		// 		Host:   config.Config.VectorstoreConfig.BaseURL,
		// 		Port:   6664,
		// 		APIKey: os.Getenv("QDRANT_API_KEY"),
		// 	})

		// 	_ = embeddingsResp
		// 	// fmt.Printf("embeddings: %v\n", embeddingsResp)

		// 	return nil
		// },
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if logger.Log != nil {
			logger.Log.Error("An error occurred", "error", err)
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.toml", "path to config.toml file")
	rootCmd.PersistentFlags().StringVarP(&logFile, "log", "l", "logs/ragtime.log", "path to log file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "include more details in the logs")

	rootCmd.AddCommand(ingest.IngestCmd)
	rootCmd.AddCommand(chat.ChatCmd)
}

func initLogs() error {
	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	logger.New(file, verbose)

	return nil
}

func initConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	// Load config file
	if cfgFile != "" {
		_, err := toml.DecodeFile(cfgFile, &config.Config)
		if err != nil {
			return err
		}

		if err := config.Config.Validate(); err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("No config file specified.")
}
