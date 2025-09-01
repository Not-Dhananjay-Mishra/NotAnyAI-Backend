package codingmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"server/models"
	"server/tools"
	"server/utils"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

type PostCodeResponse struct {
	Components map[string]string `json:"components"`
}

const syspromptpost = `
You are a coding assistant. Fix the given React code and return all components in a JSON object.

Rules:
1. Keys = filenames: "App.js" for main, others end with ".jsx".
2. Values = full valid React component code.
3. Only use React + Tailwind (className) + react-dom 18.2.0.
4. No external libs, no comments, no markdown, no extra text.
5. Tailwind classes must be valid and visually consistent (modern, cool, vibing together).
6. JSON format:

{
  "App.js": "<code>",
  "Component.jsx": "<code>"
}
also make sure that all the components are interlinked and navigation is done properly
7. No backticks. Output ONLY JSON, inside tools.
8. Navigation must be handled with React state or react-dom 18.2.0 and conditional rendering instead of react-router-dom.
9. make the home page very attractive and modern and full of content.
`

func MapToContent(m map[string]string, conn *websocket.Conn) []*genai.Content {
	var ans []*genai.Content
	for _, i := range m {
		parts := []*genai.Part{
			genai.NewPartFromText(i),
		}
		sus := genai.NewContentFromParts(parts, genai.RoleUser)
		ans = append(ans, sus)
	}
	return ans
}

func CodingPostProcessor(content []*genai.Content, conn *websocket.Conn, prompt string, allFiles []string, rag string) map[string]string {
	// Send initial processing message
	if err := conn.WriteJSON(map[string]string{"processing": "⏳ Final processing, this may take a few minutes..."}); err != nil {
		fmt.Println("WebSocket write error:", err)
	}
	time.Sleep(time.Second * 100)

	c := models.GeminiModel()
	ctx := context.Background()
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{FunctionDeclarations: []*genai.FunctionDeclaration{&tools.PostCode}},
		},
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: syspromptpost}},
		},
	}

	result, err := c.Models.GenerateContent(ctx, "gemini-2.5-flash", content, config)
	if err != nil {
		fmt.Println("GenerateContent error:", utils.Red(err))
		return nil
	}

	//temp, _ := json.Marshal(result)
	//fmt.Println("PostProcessor Raw Response:", string(temp))
	if len(result.Candidates) == 0 ||
		result.Candidates[0].Content == nil ||
		len(result.Candidates[0].Content.Parts) == 0 ||
		result.Candidates[0].Content.Parts[0].FunctionCall == nil {
		fmt.Println("Unexpected empty result from Gemini")
		sus, _ := NotCodeGen(prompt, allFiles, "empty rag", conn)
		return sus
	}

	// Extract args
	res, err := json.Marshal(result.Candidates[0].Content.Parts[0].FunctionCall.Args)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return nil
	}

	var response PostCodeResponse
	if err := json.Unmarshal(res, &response); err != nil {
		fmt.Println("Unmarshal error:", err)
		return nil
	}

	// Token usage log
	if result.UsageMetadata != nil {
		fmt.Println(utils.Red("Post Token - ", result.UsageMetadata.TotalTokenCount))
	}

	// Final message
	if err := conn.WriteJSON(map[string]string{"processing": "🎉 Done!"}); err != nil {
		fmt.Println("WebSocket write error:", err)
	}

	return response.Components
}
