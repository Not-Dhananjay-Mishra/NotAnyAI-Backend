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
	FrontendCode map[string]string `json:"frontendCode"`
	BackendCode  map[string]string `json:"backendCode"`
}

const syspromptpost = `
You are a coding assistant. Fix the given NextJS/React code and return all components in a JSON object matching this schema:

Rules:
1. Frontend pages go under "frontendCode". Keys = filenames like App.js, index.js, shop.js. Values = full valid React/NextJS code.
2. Backend API files go under "backendCode". Keys = filenames like hello.js, cart.js. Values = full valid API code (assume they are inside pages/api/).
4. Only use React + Tailwind (className) + react-dom 18.2.0 + framer-motion 11.2.6. No other libraries.
5. No comments, markdown, extra text, or explanations.
6. Tailwind classes must be consistent, modern, visually appealing.
7. All components must be interlinked and navigation must be handled via React state or conditional rendering. Do not use react-router-dom.
8. All components must be functional components with no errors.
9. Output JSON strictly in the format:

{
  "frontendCode": {
    "App.js": "<code>",
    "index.js": "<code>"
  },
  "backendCode": {
    "hello.js": "<code>"
  },
}

10. Do not use backticks, semicolons, or contractions in the text. Output ONLY JSON, inside tools.
11. Always output plain code without backticks, markdown fences, string interpolation markers, or extra escape characters.
12. all code must be interconnected like api (backend) and pages (frontend) and code should work properly as code directly go to nextjs
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

func CodingPostProcessor(content []*genai.Content, conn *websocket.Conn, prompt string, allFiles []string, rag string) PostCodeResponse {
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
		sus, _ := NotCodeGen(prompt, allFiles, "empty rag", conn)
		return sus
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
		return PostCodeResponse{}
	}

	var response PostCodeResponse
	if err := json.Unmarshal(res, &response); err != nil {
		fmt.Println("Unmarshal error:", err)
		return PostCodeResponse{}
	}

	// Token usage log
	if result.UsageMetadata != nil {
		fmt.Println(utils.Red("Post Token - ", result.UsageMetadata.TotalTokenCount))
	}

	// Final message
	if err := conn.WriteJSON(map[string]string{"processing": "🎉 Done!"}); err != nil {
		fmt.Println("WebSocket write error:", err)
	}
	PrintCodeDebug(response)
	return response
}

func PrintCodeDebug(data PostCodeResponse) {
	fmt.Println("Frontend Code:")
	for filename := range data.FrontendCode {
		fmt.Printf("File: %s\n", filename)
	}
	fmt.Println("Backend Code:")
	for filename := range data.BackendCode {
		fmt.Printf("File: %s\n", filename)
	}
}
