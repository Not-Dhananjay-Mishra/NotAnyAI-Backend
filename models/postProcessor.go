package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/tools"
	"server/utils"

	"google.golang.org/genai"
)

func PostProcessing(c *genai.Client, username string, content []*genai.Content, lastquery string, prompt []*genai.Content) string {
	ctx := context.Background()
	//content = append(content, genai.NewContentFromText(lastquery, genai.RoleUser))
	result, err := c.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		content,
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				{
					FunctionDeclarations: []*genai.FunctionDeclaration{
						&tools.ToolDeciderAgent},
				},
			},
			SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: "These are response from fucntion calls post process it to give final output to user, if more info needed u can call tool but call only if u need more infomation"}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	tkn, _ := json.Marshal(result.UsageMetadata.TotalTokenCount)
	fmt.Println("Total Token used: ", string(tkn))
	if result.Candidates[0].Content.Parts[0].Text != "" {
		res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].Text)
		utils.MemoryStore[username] = append(utils.MemoryStore[username], genai.NewContentFromText(string(res), genai.RoleModel))
		fmt.Println(utils.Yellow("AI : "), string(res))
		return string(res)
	} else {
		res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].FunctionCall.Args)
		var data Agent
		json.Unmarshal(res, &data)
		fmt.Println(utils.Cyan(string(res)))
		content := ToolCaller(data, prompt[len(prompt)-1].Parts[0].Text)
		sus := PostProcessing(c, username, content, prompt[len(prompt)-1].Parts[0].Text, prompt)
		return sus
	}
}
