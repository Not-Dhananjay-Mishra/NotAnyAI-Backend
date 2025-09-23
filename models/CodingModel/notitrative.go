package codingmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/models"
	"server/tools"
	"server/utils"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

type GeneratedCode struct {
	FrontendCode map[string]string `json:"frontendCode"`
	BackendCode  map[string]string `json:"backendCode"`
	Components   map[string]string `json:"components"`
}

func NotCodeGen(prompt string, allFiles []string, rag string, conn *websocket.Conn) (PostCodeResponse, int32) {
	if err := conn.WriteJSON(map[string]string{"processing": "Model Overloaded retrying final time please wait..."}); err != nil {
		fmt.Println("WebSocket write error:", err)
	}
	time.Sleep(time.Second * 60)
	// Build system prompt
	sysprompt := fmt.Sprintf(`
You are a Next.js code assistant.

The project has:
- Frontend React .js pages and Backend API route files: %v

Generate ONLY the full, valid code.

Rules:

Frontend:
- Keys = filenames inside "pages/" (e.g., "index.js", "about.jsx").
- Values = full valid React component code.
- Use React + Tailwind (className) + react-dom 18.2.0 only.
- No external libs, no comments, no markdown, no extra text.
- Tailwind classes must be valid and visually consistent (modern, cool, vibing together).
- Navigation must be handled with React state or conditional rendering instead of react-router-dom.
- Apply RAG intelligently: adjust Tailwind colors for a modern/cool theme, not copy-paste. If you don't want RAG: %v.
- Use img URLs only if needed, styled with proper Tailwind given in RAG.
- Try to add animations and interactivity to make the UI more engaging.
- Use three.js with react-three-fiber and drei for 3D models and animations if needed.
- Add SVGs for icons and illustrations.
- Do not use comments in the code.

Backend:
- Keys = filenames inside "pages/api/" (e.g., "hello.js").
- Values = full valid Next.js API route code (export default handler).
- Must return JSON responses.
- Use only Node.js built-in modules and Next.js conventions.
- No external libraries beyond those already specified.

JSON format:

{
  "frontendCode": {
    "index.js": "<code>",
    "about.jsx": "<code>"
  },
  "backendCode": {
    "hello.js": "<code>",
    "auth.js": "<code>"
  }
}

No backticks. Output ONLY JSON, inside tools.
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
		return PostCodeResponse{}, 0
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
			return PostCodeResponse{}, int32(result.UsageMetadata.TotalTokenCount)
		}
		return PostCodeResponse{}, 0
	}

	part := result.Candidates[0].Content.Parts[0]
	if part.FunctionCall == nil || part.FunctionCall.Args == nil {
		res, _ := json.Marshal(part)
		fmt.Println("FunctionCall or Args nil. Part:", string(res))
		if result.UsageMetadata != nil {
			return PostCodeResponse{}, int32(result.UsageMetadata.TotalTokenCount)
		}
		return PostCodeResponse{}, 0
	}

	// Marshal the function args to JSON and unmarshal into GeneratedCode
	argsJSON, _ := json.Marshal(part.FunctionCall.Args)
	var gen PostCodeResponse
	if err := json.Unmarshal(argsJSON, &gen); err != nil {
		// If direct unmarshal to struct fails, try to inspect the raw args
		fmt.Println("Failed to unmarshal args to GeneratedCode:", err)
		fmt.Println("Raw args:", string(argsJSON))
		if result.UsageMetadata != nil {
			return PostCodeResponse{}, int32(result.UsageMetadata.TotalTokenCount)
		}
		return PostCodeResponse{}, 0
	}

	tkn := int32(0)
	if result.UsageMetadata != nil {
		tkn = int32(result.UsageMetadata.TotalTokenCount)
	}
	if err := conn.WriteJSON(map[string]string{"processing": "🎉 Done!"}); err != nil {
		fmt.Println("WebSocket write error:", err)
	}
	return gen, tkn
}
