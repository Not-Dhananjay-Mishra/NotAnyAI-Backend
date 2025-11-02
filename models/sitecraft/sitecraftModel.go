package sitecraft

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/utils"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/gorilla/websocket"
)

type AIPostCodeResponse struct {
	FrontendCode map[string]string `json:"frontendCode"`
	BackendCode  map[string]string `json:"backendCode"`
}

// Struct representing an incoming request to generate code
type Sus struct {
	Query string `json:"query"`
}

// Struct for the AI model’s raw output
type GenAIResponse struct {
	FrontendFile []string `json:"frontendFiles"`
	BackendFile  []string `json:"backendFiles"`
}

// Struct for the formatted response we’ll send to the frontend
type FileData struct {
	Filename string `json:"filename" jsonschema_description:"Name of the code file dont add pages/ or api/ just the filename with .js extension and dont make nested files(like product/[id].js) only simple filenames"`
	Content  string `json:"content" jsonschema_description:"File content "`
}

type PostCodeResponse struct {
	FrontendCode []FileData `json:"frontendCode" jsonschema_description:"List of frontend files"`
	BackendCode  []FileData `json:"backendCode" jsonschema_description:"List of backend files"`
}

// Struct for handling in-progress jobs
type Process struct {
	Data    Sus
	ReplyCh chan AIPostCodeResponse
	Conn    *websocket.Conn
}

// Buffered channel for pending code-generation requests
var UnderProcessCodeNew chan Process

// Initialize queue
func init() {
	UnderProcessCodeNew = make(chan Process, 50)
}

// Processor continuously listens for code generation requests
func Processor() {
	for {
		data := <-UnderProcessCodeNew
		//data.Conn.WriteJSON(utils.Response{Text: "Processing your code generation request..."})
		fmt.Println(data.Data.Query)
		if data.Data.Query != "" {
			fmt.Println(utils.Blue(data.Data.Query))
			sus := SiteCraftAIModel(data.Data, data.Conn)
			data.ReplyCh <- sus
		}
	}
}

// Core function that interacts with Google’s Gemini via Genkit
func SiteCraftAIModel(data Sus, conn *websocket.Conn) AIPostCodeResponse {
	ctx := context.Background()
	fmt.Println("Generating code for query:", data.Query)
	// Initialize Genkit with Google Gemini
	g := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.5-pro"), // reliable full model
	)

	prompt := fmt.Sprintf(`
You are a frontend and backend code assistant for generating Next.js projects.
For any given prompt, determine only the necessary files required to implement the described website.
- Frontend pages must be returned in "frontendFiles" as .js files (e.g., index.js, about.js).
- Backend API endpoints must be returned in "backendFiles" as .js files under (e.g., hello.js).
- Do not generate code, explanations, or extra text—only return valid file names via the provided tool.
- Keep the count of files to the absolute minimum required to fulfill the user's request.
- dont create nested file execpt pages and api only <filename>.js and <filename>.js nothing other than that
- No follow-up questions.
User request: "%s"

Generate:
1. Frontend code (NextJS + Tailwind) as JS files.
2. Backend code (NextJS + JSON API).
Return JSON in this exact format:
Allowed libraries and environment:
2) Frontend(Next.js): React 18.2.0, react-dom 18.2.0, framer-motion 11.2.6, Tailwind CSS. No other libraries.
3) Backend: Next.js API route handlers only. No external packages.

Frontend rules:
4) Use functional React components with a default export for pages.
5) Use Tailwind utilities for styling; avoid inline styles.
6) Use inline SVGs for icons; do not import icon libraries.
7) Use useState and useEffect only for state; keep state minimal.
8) Include responsive layouts, keyboard accessibility, semantic HTML, aria labels, and visible focus states.
9) Provide placeholder loading states for async data.
10) Add subtle animations and micro-interactions using framer-motion and Tailwind transitions.
11) Each file must implement a single responsibility: one top-level component or one API handler.

Backend rules:
12) If the filename suggests an API route (for example path contains /api/ or filename ends with .js intended for /api), produce a valid Next.js API handler that:
   - exports a function compatible with Next's API routes,
   - validates inputs, handles errors, returns JSON responses,
   - uses no external packages.
{
  "frontendCode": {
    "App.js": "<code>",
    "index.js": "<code>"
  },
  "backendCode": {
    "hello.js": "<code>"
  },
}
`, data.Query)
	var typee PostCodeResponse
	resp, err := genkit.Generate(ctx, g, ai.WithPrompt(prompt), ai.WithOutputType(&typee))
	if err != nil {
		log.Printf("Error generating code: %v", err)
		sendError(conn, "AI model failed to generate code.")
		return AIPostCodeResponse{}
	}
	json.Unmarshal([]byte(resp.Text()), &typee)
	var frontendCodeMap = make(map[string]string)
	for _, file := range typee.FrontendCode {
		frontendCodeMap[file.Filename] = file.Content
	}
	var backendCodeMap = make(map[string]string)
	for _, file := range typee.BackendCode {
		backendCodeMap[file.Filename] = file.Content
	}
	return AIPostCodeResponse{
		FrontendCode: frontendCodeMap,
		BackendCode:  backendCodeMap,
	}
}

// Sends data to the WebSocket connection safely
func sendWSMessage(conn *websocket.Conn, msg interface{}) {
	if conn == nil {
		return
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling WS message: %v", err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, bytes)
}

// Sends an error message to WebSocket
func sendError(conn *websocket.Conn, msg string) {
	if conn == nil {
		return
	}
	errorMsg := map[string]string{"error": msg}
	bytes, _ := json.Marshal(errorMsg)
	conn.WriteMessage(websocket.TextMessage, bytes)
}
