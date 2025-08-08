package models

import (
	"context"
	"log"
	"server/utils"

	"google.golang.org/genai"
)

func GeminiModel() *genai.Client {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  utils.GEMINI_API,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}
	return client
}
