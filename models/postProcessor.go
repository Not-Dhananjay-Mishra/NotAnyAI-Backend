package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/tools"
	"server/utils"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

func PostProcessing(c *genai.Client, username string, content []*genai.Content, lastquery string, prompt []*genai.Content, conn *websocket.Conn) string {
	ctx := context.Background()
	conn.WriteJSON(utils.Response{Text: "Processing request..."})
	sus := `You have received responses from multiple function calls related to the user’s query. Your task is to:
	- Combine and process these responses to generate a clear, concise, and relevant final answer for the user with all infomation needed.
	- Only and only request additional information or call other tools if the current responses are insufficient to fully answer the query avoid if u can dont call same tool with same query again.
	- Avoid unnecessary repetition or unrelated details.
	- Present the final output in a user-friendly and informative way.`
	//content = append(content, genai.NewContentFromText(lastquery, genai.RoleUser))
	result, err := c.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		content,
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				{
					FunctionDeclarations: []*genai.FunctionDeclaration{
						&tools.ToolDeciderAgent},
				},
			},
			SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: sus}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	tkn, _ := json.Marshal(result.UsageMetadata.TotalTokenCount)
	fmt.Println("Total Token used: ", string(tkn))
	part := result.Candidates[0].Content.Parts[0]

	if part.Text != "" {
		utils.MemoryStore[username] = append(utils.MemoryStore[username], genai.NewContentFromText(part.Text, genai.RoleModel))
		fmt.Println(utils.Yellow("AI : "), part.Text)
		return part.Text
	} else {
		res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].FunctionCall.Args)
		var data Agent
		json.Unmarshal(res, &data)
		fmt.Println(utils.Cyan(string(res)))
		content := ToolCaller(data, prompt[len(prompt)-1].Parts[0].Text, conn)
		sus := PostProcessing(c, username, content, prompt[len(prompt)-1].Parts[0].Text, prompt, conn)
		return sus
	}
}
