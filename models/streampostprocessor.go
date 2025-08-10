package models

import (
	"context"
	"fmt"
	"log"
	"server/utils"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

func StreamPostProcessing(c *genai.Client, username string, content []*genai.Content, lastquery string, prompt []*genai.Content, conn *websocket.Conn) string {
	log.Println("Reached SPP")
	ctx := context.Background()
	conn.WriteJSON(utils.Response{Text: "Processing request..."})
	sus := `You have received responses from multiple function calls related to the user’s query. Your task is to:
	- Combine and process these responses to generate a clear, concise, and relevant final answer for the user with all infomation needed.
	- Only and only request additional information or call other tools if the current responses are insufficient to fully answer the query avoid if u can dont call same tool with same query again.
	- Avoid unnecessary repetition or unrelated details.
	- Present the final output in a user-friendly and informative way.`
	//content = append(content, genai.NewContentFromText(lastquery, genai.RoleUser))
	fmt.Println(content)
	result := c.Models.GenerateContentStream(
		ctx,
		"gemini-2.5-flash",
		content,
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				{
					GoogleSearch: &genai.GoogleSearch{},
				},
			},
			SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: sus}}},
		},
	)

	/*tkn, _ := json.Marshal(result.UsageMetadata.TotalTokenCount)
	fmt.Println("Total Token used: ", string(tkn))
	part := result.Candidates[0].Content.Parts[0]*/
	finalans := ""
	for chunk, _ := range result {
		//res, _ := json.Marshal(chunk)
		//log.Println(string(res))
		part := chunk.Candidates[0].Content.Parts[0]
		if part.Text != "" {
			//fmt.Println(utils.Yellow("AI : "), part.Text)
			finalans += part.Text
			conn.WriteJSON(utils.Response{Text: finalans})
			//return part.Text
		}
	}
	utils.MemoryStore[username] = append(utils.MemoryStore[username], genai.NewContentFromText(finalans, genai.RoleModel))
	return finalans
}
