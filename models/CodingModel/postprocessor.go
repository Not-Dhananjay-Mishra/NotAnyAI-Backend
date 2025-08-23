package codingmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"server/models"
	"server/tools"
	"server/utils"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

type PostCodeResponse struct {
	Components map[string]string `json:"components"`
}

const syspromptpost = `
You are a coding assistant.
Your task is to fix the given React code and return every React component in a JSON object.

### Rules:
1. Each key must be the filename:
   - Use "App.js" for the main entry file (only this file should end with .js).
   - Use ".jsx" for all other components.
2. Each value must be the **full valid React component code** for that file.
3. Do not use any external libraries. Only use:
   - React
   - Tailwind CSS (via className strings)
4. Do not output explanations, comments, or Markdown formatting.
5. All Tailwind CSS classes must be valid.
6. The JSON must follow this structure exactly:

{
  "App.js": "<entire React code here>",
  "ComponentName.jsx": "<entire React code here>",
  ...
}

7. Do not wrap the code in triple backticks.
8. Do not include extra text outside the JSON object.
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

func CodingPostProcessor(content []*genai.Content, conn *websocket.Conn) map[string]string {
	conn.WriteJSON(utils.Response{Text: "⏳ Final processing, this may take a few minutes..."})
	c := models.GeminiModel()
	ctx := context.Background()
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{FunctionDeclarations: []*genai.FunctionDeclaration{&tools.PostCode}},
		},
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: syspromptpost}}},
	}
	fmt.Println(content)
	result, err := c.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		content,
		config,
	)
	if err != nil {
		fmt.Println(utils.Red(err))
	}
	//eee, _ := json.Marshal(result)
	//fmt.Println(utils.Magenta(string(eee)))
	res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].FunctionCall.Args)
	var response PostCodeResponse
	json.Unmarshal(res, &response)
	tkn := result.UsageMetadata.TotalTokenCount
	fmt.Println(utils.Red("Post Token - ", tkn))
	conn.WriteJSON(utils.Response{Text: "🎉 Done!"})
	return response.Components
}
