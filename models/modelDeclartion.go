package models

import (
	"context"
	"encoding/json"
	"fmt"
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

func PlainLLM(c *genai.Client, prompt []*genai.Content, username string) string {
	ctx := context.Background()
	sus := `answer the query, give full response`
	//content = append(content, genai.NewContentFromText(lastquery, genai.RoleUser))
	result, err := c.Models.GenerateContent(
		ctx,
		"gemini-2.5-pro",
		prompt,
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: sus}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		res, _ := json.Marshal(result)
		log.Println(string(res))
		for i := range prompt {
			log.Println(prompt[i].Parts)
		}
		log.Println("No candidates or parts returned from model")
		return "Sorry, I couldn't find any information for that."
	}

	part := result.Candidates[0].Content.Parts[0]

	if part.Text != "" {
		utils.MemoryStore[username] = append(utils.MemoryStore[username], genai.NewContentFromText(part.Text, genai.RoleModel))
		fmt.Println(utils.Yellow("AI : "), part.Text)
		return part.Text
	}
	return ""
}
