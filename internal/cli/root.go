package cli

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/abspayd/ragtime/internal/config"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	cfgFile string

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
			return initConfig()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.toml", "Path to config.toml file")
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
