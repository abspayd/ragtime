package documents

import (
	"context"
	"os"
	"os/signal"

	"github.com/abspayd/ragtime/internal/logger"
)

func UploadDocuments(paths []string) error {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// context.Background

	for _, path := range paths {

		data, err := Load(ctx, path)
		if err != nil {
			return err
		}

		logger.Log.Info("read from file", path, data)
	}

	return nil
}
