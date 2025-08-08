package models

import (
	"context"
	"encoding/json"
	"fmt"
	"server/utils"

	"google.golang.org/genai"
)

func CacheContext(client *genai.Client, prompt string) {
	ctx := context.Background()

	modelName := "gemini-2.0-flash"

	var memory []*genai.Content

	memory = append(memory, genai.NewContentFromText("my name is xyz", genai.RoleUser))
	memory = append(memory, genai.NewContentFromText("ok .....", genai.RoleModel))
	memory = append(memory, genai.NewContentFromText("what is my name", genai.RoleUser))

	response, _ := client.Models.GenerateContent(
		ctx,
		modelName,
		memory,
		&genai.GenerateContentConfig{
			SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: "You are a helpful assistant."}}},
		},
	)
	tkn, _ := json.Marshal(response.UsageMetadata.TotalTokenCount)
	fmt.Println("Total Token used: ", string(tkn))
	res, _ := json.Marshal(response.Candidates[0].Content.Parts[0].Text)
	fmt.Println(utils.Yellow("AI : "), string(res))
}
func AddToMemoryUSER(username string, prompt string) {
	utils.MemoryStore[username] = append(utils.MemoryStore[username], genai.NewContentFromText(prompt, genai.RoleUser))
}
