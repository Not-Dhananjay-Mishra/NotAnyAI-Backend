package codingmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/models"
	"server/tools"
	"server/utils"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

type GeneratedCode struct {
	Code map[string]string `json:"components"`
}

func NotCodeGen(prompt string, allFiles []string, rag string, conn *websocket.Conn) (map[string]string, int32) {
	if err := conn.WriteJSON(map[string]string{"processing": "Model Overloaded retrying final time please wait..."}); err != nil {
		fmt.Println("WebSocket write error:", err)
	}
	// Build system prompt
	sysprompt := fmt.Sprintf(`
You are a frontend code assistant.

The project has React .jsx files: %v. Generate ONLY the full, valid code.

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

7. No backticks. Output ONLY JSON, inside tools.
8. Navigation must be handled with React state or react-dom 18.2.0 and conditional rendering instead of react-router-dom.
Apply RAG intelligently: adjust Tailwind colors for a modern/cool theme, not copy-paste. RAG: %v.
Use img URLs only if needed, styled with proper Tailwind given in RAG.
Try to add animations and interactivity to make the UI more engaging.
Use three.js three-js-react for 3d models and animations if needed.
Add svg for icons and illustrations.
Do not use comments in the code.
Use images (make your own svg and add that) and animations to make the UI more engaging.
`, allFiles, rag)

	// Create client
	c := models.GeminiModel()
	ctx := context.Background()

	// system instruction
	system := &genai.Content{
		Parts: []*genai.Part{
			{Text: sysprompt},
		},
	}

	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{FunctionDeclarations: []*genai.FunctionDeclaration{&tools.PostCode}},
		},
		SystemInstruction: system,
		// You can set other config fields (temperature, max tokens) here if supported
	}

	// Call the API
	result, err := c.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), config)
	if err != nil {
		log.Println("GenerateContent error:", utils.Red(err))
		conn.WriteJSON(map[string]string{"processing": "Error from model, please try again later."})
		return nil, 0
	}

	// Debug print of raw response (safe for local debugging)
	temp, _ := json.Marshal(result)
	fmt.Println("PostProcessor Raw Response:", string(temp))

	// Basic validations
	if len(result.Candidates) == 0 ||
		result.Candidates[0].Content == nil ||
		len(result.Candidates[0].Content.Parts) == 0 ||
		result.Candidates[0].Content.Parts[0].FunctionCall == nil {
		log.Println("Unexpected empty result from Gemini or missing function call")
		// Try to return token usage if available
		if result.UsageMetadata != nil {
			return nil, int32(result.UsageMetadata.TotalTokenCount)
		}
		return nil, 0
	}

	part := result.Candidates[0].Content.Parts[0]
	if part.FunctionCall == nil || part.FunctionCall.Args == nil {
		res, _ := json.Marshal(part)
		fmt.Println("FunctionCall or Args nil. Part:", string(res))
		if result.UsageMetadata != nil {
			return nil, int32(result.UsageMetadata.TotalTokenCount)
		}
		return nil, 0
	}

	// Marshal the function args to JSON and unmarshal into GeneratedCode
	argsJSON, _ := json.Marshal(part.FunctionCall.Args)
	var gen GeneratedCode
	if err := json.Unmarshal(argsJSON, &gen); err != nil {
		// If direct unmarshal to struct fails, try to inspect the raw args
		fmt.Println("Failed to unmarshal args to GeneratedCode:", err)
		fmt.Println("Raw args:", string(argsJSON))
		if result.UsageMetadata != nil {
			return nil, int32(result.UsageMetadata.TotalTokenCount)
		}
		return nil, 0
	}

	tkn := int32(0)
	if result.UsageMetadata != nil {
		tkn = int32(result.UsageMetadata.TotalTokenCount)
	}
	if err := conn.WriteJSON(map[string]string{"processing": "🎉 Done!"}); err != nil {
		fmt.Println("WebSocket write error:", err)
	}
	return gen.Code, tkn
}
