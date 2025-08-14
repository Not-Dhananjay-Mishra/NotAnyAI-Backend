package models

import (
	"context"
	"log"
	"server/utils"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

func ImgStreamPostProcessing(c *genai.Client, username string, content []*genai.Content, lastquery string, prompt []*genai.Content, conn *websocket.Conn) string {
	log.Println("Reached SPP ")
	ctx := context.Background()
	conn.WriteJSON(utils.Response{Text: "Processing request..."})
	sus := `You may receive results from multiple function calls related to the same user request, or a direct user request. Your task is to:
		Combine & Process: Merge the outputs from all available responses into a single, coherent, and relevant final result.
		Enhance Clarity: Ensure the processed image output or related information is presented in a clear, concise, and user-friendly way.
		Avoid Redundancy: Do not repeat the same tool call for the same query, and avoid duplicating similar information.
		Completeness First: Only request additional data or trigger new processing if the existing responses are insufficient to fully meet the user’s needs.
		Direct Handling: If a direct query is given and can be answered with existing knowledge, handle it without extra tool calls.
		User Presentation: Deliver the final processed image or related explanation in a clean, understandable manner without unrelated details.`
	//content = append(content, genai.NewContentFromText(lastquery, genai.RoleUser))
	//log.Println(content)
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
		if chunk.Candidates[0].Content.Parts[0].Text != "" {
			//fmt.Println(utils.Yellow("AI : "), part.Text)
			finalans += chunk.Candidates[0].Content.Parts[0].Text
			conn.WriteJSON(utils.Response{Text: finalans})
			//return part.Text
		} else {
			conn.WriteJSON(utils.Response{Text: "error occur refresh page"})
		}
	}
	utils.MemoryStore[username] = append(utils.MemoryStore[username], genai.NewContentFromText(finalans, genai.RoleModel))
	log.Println("Done SPP ")
	return finalans
}
