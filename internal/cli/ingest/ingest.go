package ingest

import (
	"github.com/abspayd/ragtime/internal/documents"
	"github.com/spf13/cobra"
)

var (
	noGitIgnore bool
	IgnoreGlob  string

	IngestCmd = &cobra.Command{
		Use:   "ingest path...",
		Short: "Add documents to the vector store",
		Long:  "Add documents for your LLMs to use as a reference",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := documents.UploadDocuments(args); err != nil {
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
