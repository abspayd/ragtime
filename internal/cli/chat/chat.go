package chat

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/abspayd/ragtime/internal/config"
	"github.com/abspayd/ragtime/internal/models"
	"github.com/qdrant/go-client/qdrant"
	"github.com/spf13/cobra"
)

var (
	ChatCmd = &cobra.Command{
		Use:   "chat message...",
		Short: "Chat with an LLM using RAG search",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			message := strings.Join(args, " ")

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
			defer cancel()

			chatClient := models.NewOpenAIClient(config.Config.ChatModelConfig.Model, config.Config.ChatModelConfig.BaseURL, os.Getenv("OPENAI_API_KEY"))

			embeddingsClient := models.NewOpenAIClient(config.Config.EmbeddingConfig.Model, config.Config.EmbeddingConfig.BaseURL, os.Getenv("OPENAI_API_KEY"))

			qdrantClient, err := qdrant.NewClient(&qdrant.Config{
				Host:   config.Config.VectorstoreConfig.BaseURL,
				Port:   6334,
				APIKey: os.Getenv("QDRANT_API_KEY"),
				UseTLS: false,
			})
			if err != nil {
				return err
			}

			embeddings, err := embeddingsClient.Embed(ctx, message)

			rag_results, err := qdrantClient.Query(ctx, &qdrant.QueryPoints{
				CollectionName: config.Config.VectorstoreConfig.Collection,
				Query:          qdrant.NewQuery(embeddings[0]...),
				WithPayload:    qdrant.NewWithPayload(true),
			})
			if err != nil {
				return err
			}

			var buf bytes.Buffer

			buf.WriteString("<Context>")
			for _, result := range rag_results {
				if result.GetScore() < 0.5 {
					continue
				}

				payload := result.GetPayload()

				fmt.Fprintf(&buf, "<Path>%s</Path>", payload["path"])
				fmt.Fprintf(&buf, "<Score>%f</Score>", result.GetScore())
				fmt.Fprintf(&buf, "<Text>%s</Text>", payload["text"])
			}
			buf.WriteString("</Context>")
			rag_context := buf.String()

			system_prompt := "You are an assistant who uses only the information provided to you as your source of truth." +
				"Each piece of context will include the following:" +
				"  - Path: the file path where the information was collected" +
				"  - Score: the confidence level of how relevant the data is to the user's query (scored from 0 to 1)" +
				"  - Text: the fragment of text containing the context data" +
				"You always double check your work, and make it clear when you don't know." +
				"You must include a sitation of the file path name if you use its information in your response." +
				"Respond in a highly factual manner and avoid using emojis."

			response, err := chatClient.Chat(ctx, []models.Message{
				{
					Role:    "system",
					Content: system_prompt + rag_context,
				},
				{
					Role:    "user",
					Content: message,
				},
			})
			if err != nil {
				return err
			}

			fmt.Println(response)

			return nil
		},
	}
)

func init() {

}
