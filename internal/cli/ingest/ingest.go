package ingest

import (
	"github.com/abspayd/ragtime/internal/config"
	"github.com/abspayd/ragtime/internal/documents"
	"github.com/abspayd/ragtime/internal/logger"
	"github.com/abspayd/ragtime/internal/models"
	"github.com/qdrant/go-client/qdrant"
	"github.com/spf13/cobra"
)

var (
	noGitIgnore bool
	IgnoreGlob  string

	IngestCmd = &cobra.Command{
		Use:   "ingest path...",
		Short: "Add documents to the vector store",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			embeddingsClient := models.NewOpenAIClient(config.Config.EmbeddingConfig.Model, config.Config.EmbeddingConfig.BaseURL, "")

			qdrantClient, err := qdrant.NewClient(&qdrant.Config{
				Host:   config.Config.VectorstoreConfig.BaseURL,
				Port:   6334,
				APIKey: "",
				UseTLS: false,
			})
			if err != nil {
				return err
			}

			logger.Log.Info("connected to qdrant")

			if err := documents.UploadDocuments(args, config.Config.VectorstoreConfig.Collection, embeddingsClient, qdrantClient); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	IngestCmd.PersistentFlags().BoolVar(&noGitIgnore, "no-git-ignore", false, "do not ignore files in '.gitignore'")
	IngestCmd.PersistentFlags().StringVarP(&IgnoreGlob, "ignore", "I", "", "glob pattern to exclude files or directories")
}
