package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/BurntSushi/toml"
	"github.com/abspayd/ragtime/internal/config"
	"github.com/abspayd/ragtime/internal/logger"
	"github.com/abspayd/ragtime/internal/models"
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
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			client := models.NewOpenAIClient("gpt-oss", config.Config.ChatModelConfig.BaseURL, "")
			resp, err := client.Chat(ctx, []models.Message{
				{
					Role:    "user",
					Content: "Hello",
				},
			})

			if err != nil {
				return err
			}

			fmt.Println(resp)
			for _, choice := range resp.Choices {
				fmt.Printf("%s: %s\n", choice.Message.Role, choice.Message.Content)
			}

			return nil
		},
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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.toml", "Path to config.toml file")
	rootCmd.PersistentFlags().StringVarP(&logFile, "log", "l", "logs/ragtime.log", "Path to log file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Include more details in the logs")
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
