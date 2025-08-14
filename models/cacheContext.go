package models

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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
	for len(utils.MemoryStore[username]) > 20 {
		utils.MemoryStore[username] = utils.MemoryStore[username][1:]
	}
}
func AddImgToMemoryUSER(username string, prompt string, path string) {
	bytes, _ := os.ReadFile(path)
	parts := []*genai.Part{
		genai.NewPartFromBytes(bytes, "image/png"),
		genai.NewPartFromText(prompt),
	}

	utils.MemoryStore[username] = append(utils.MemoryStore[username], genai.NewContentFromParts(parts, genai.RoleUser))
	//fmt.Println(utils.MemoryStore[username])
	for len(utils.MemoryStore[username]) > 20 {
		utils.MemoryStore[username] = utils.MemoryStore[username][1:]
	}
}

func PrintMemobyUsername(username string) {
	mem := utils.MemoryStore[username]
	for _, i := range mem {
		if i.Parts[0].Text != "" {
			fmt.Println(i.Parts[0].Text)
		} else if i.Parts[1].Text != "" {
			fmt.Println(i.Parts[1].Text)
		} else if i.Parts[0].FunctionCall.Name != "" {
			fmt.Println(i.Parts[0].FunctionCall.Name)
		} else if i.Parts[0].FunctionResponse.Name != "" {
			res, _ := json.Marshal(i.Parts[0].FunctionResponse)
			fmt.Println(string(res))
		}
		fmt.Println(i.Role)
	}
}
func PrintMemobyContent(mem []*genai.Content) {
	for _, c := range mem {
		for _, part := range c.Parts {
			switch {
			case part.Text != "":
				fmt.Println(part.Text)

			case part.FunctionCall != nil && part.FunctionCall.Name != "":
				fmt.Println("Function Call:", part.FunctionCall.Name)

			case part.FunctionResponse != nil && part.FunctionResponse.Name != "":
				res, err := json.Marshal(part.FunctionResponse)
				if err != nil {
					fmt.Println("Error marshaling function response:", err)
				} else {
					fmt.Println("Function Response:", string(res))
				}
			}
		}
		fmt.Println("Role:", c.Role)
	}
}
