package models

import (
	"context"
	"log"
	"server/utils"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

func StreamPostProcessing(c *genai.Client, username string, content []*genai.Content, lastquery string, prompt []*genai.Content, conn *websocket.Conn) string {
	//log.Println("Reached SPP ")
	ctx := context.Background()
	conn.WriteJSON(utils.Response{Text: "Processing request..."})
	sus := `You have received responses from multiple function calls related to the user’s query or direct query from user. Your task is to:
	- Combine and process these responses to generate a clear, concise, and relevant final answer for the user with all information needed.
	- Only request additional information or call other tools if the current responses are insufficient to fully answer the query. Avoid calling the same tool with the same query again.
	- Avoid unnecessary repetition or unrelated details.
	- Present the final output in a user-friendly and informative way.
	- If you receive a direct query, answer directly with your own knowledge without calling tools.
	- If the user asks "Who are you?" or "Who created you?", respond exactly with: "I am NotAnyAI, made by Dhananjay Mishra. I use the Gemini AI model."`
	//content = append(content, genai.NewContentFromText(lastquery, genai.RoleUser))
	//log.Println(content[len(content)-1].Parts[0].Text)
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
	log.Println("Done SPP ")
	return finalans
}
