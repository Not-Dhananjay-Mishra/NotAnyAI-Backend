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

type GenteratedCode struct {
	Code string `json:"code"`
}

/*func Itrative(data []string, prompt string, c *genai.Client, conn *websocket.Conn) map[string]string {
	// with go rountines
	var wg sync.WaitGroup
	mu := sync.Mutex{}
	results := make(map[string]string)
	totaltkn := 0

	for _, file := range data {
		wg.Add(1)
		go func(fname string) {
			defer wg.Done()
			mu.Lock()
			//conn.WriteJSON(utils.Response{Text: "⚙️ Started generation of " + file})
			conn.WriteJSON(map[string]string{"codegenstart": "⚙️ Started generation of " + file})
			ragout := RAGQueryDecider(prompt, conn, fname)
			time.Sleep(time.Second * 5)
			code, tkn := CodeGen(prompt, fname, data, ragout)
			//conn.WriteJSON(utils.Response{Text: "✅ Completed generation of " + file})
			conn.WriteJSON(map[string]string{"codegencomplete": "✅ Completed generation of " + file})
			results[fname] = code
			totaltkn += int(tkn)
			mu.Unlock()
		}(file)
		time.Sleep(time.Second * 15)
	}
	wg.Wait()
	fmt.Println(utils.Red("Code Gen Token - ", totaltkn))
	return results
}*/

func ItrativeWithoutGo(data []string, prompt string, c *genai.Client, conn *websocket.Conn) map[string]string {
	// without go rountines
	results := make(map[string]string)
	totaltkn := 0
	ragout := RAGQueryDecider(prompt, conn, "")

	for _, file := range data {
		conn.WriteJSON(map[string]string{"codegenstart": "⚙️ Started generation of " + file})
		time.Sleep(time.Second * 2)

		code, tkn := CodeGen(prompt, file, data, ragout)
		conn.WriteJSON(map[string]string{"codegencomplete": "✅ Completed generation of " + file})
		results[file] = code
		totaltkn += int(tkn)
		time.Sleep(time.Second * 2)
	}
	fmt.Println(utils.Red("Code Gen Token - ", totaltkn))
	return results
}

func CodeGen(prompt string, targetFile string, allFiles []string, rag string) (string, int32) {
	c := models.GeminiModel()
	ctx := context.Background()

	sysprompt := fmt.Sprintf(`You are a single-file Next.js code generator.

Inputs:
- targetFile: %s
- allFiles: %v
- ragTheme: %v

FOLLOW THIS RULE REGARDLESS OF CONDITION-
You are a precise single-file code generator: output only the exact, runnable source code for the requested target file on stdout with no extra text, 
no markdown/backticks, no surrounding JSON wrappers, no added escape characters or backslashes, do not alter or escape quotes or apostrophes, 
do not introduce contractions in generated natural-language strings, and if you cannot produce valid code return a minimal syntactically-correct file that renders a clear runtime 
error message in the UI (not as comments).

Primary goal:
1) Return only the complete, valid source code for the requested target file %s with no extra text, markup, comments, or surrounding fences. This code will be directly run 

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

RAG theme:
13) If ragTheme (%v) is empty or nil, do not force a RAG palette. If provided, adapt Tailwind colors and spacing to reflect the theme, but maintain accessibility and contrast.

Safety and produce-ready:
14) Generate production-ready code: no commented-out code, no debug console logs, no unused imports.
15) Ensure imports reference only allowed packages and that runtime references exist in the file.
16) Keep the file self-contained and ready to drop into a Next.js project.
17) VERY IMPORTANT ONLY GIVE CODE IN TOOL Response (code) not in plain text
18) Always output plain code without backticks, markdown fences, string interpolation markers, or extra escape characters.

Fallback behavior:
19) If the system cannot determine whether the target is frontend or backend, prefer producing a safe frontend page that renders a clear error/fallback UI and documents required inputs inside the UI (not as comments).
`, targetFile, allFiles, rag, targetFile, rag)

	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{FunctionDeclarations: []*genai.FunctionDeclaration{&tools.GenerateCode}},
		},
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: sysprompt}}},
	}

	result, err := c.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)

	if err != nil {
		if rag == "empty rag" {
			fmt.Println(utils.Red("Error even after retrying without RAG:", err))
			return "", 0
		}
		fmt.Println(utils.Red("Retrying without RAG due to error:", err))
		time.Sleep(time.Second * 70)
		return CodeGen(prompt, targetFile, allFiles, "empty rag")
	}

	if len(result.Candidates) == 0 {
		if rag == "empty rag" {
			fmt.Println(utils.Red("No candidates even after retrying without RAG"))
			return "", 0
		}
		time.Sleep(time.Second * 7)
		return CodeGen(prompt+" give response in tool only", targetFile, allFiles, "empty rag")
	}

	part := result.Candidates[0].Content
	if part == nil || len(part.Parts) == 0 {
		if rag == "empty rag" {
			fmt.Println(utils.Red("No content parts even after retrying without RAG"))
			return "", 0
		}
		time.Sleep(time.Second * 7)
		return CodeGen(prompt+" give response in tool only", targetFile, allFiles, "empty rag")
	}

	// Safely check function call
	fc := part.Parts[0].FunctionCall
	if fc == nil || fc.Args == nil {
		if rag == "empty rag" {
			fmt.Println(utils.Red("No valid function call args even after retrying without RAG"))
			// Print what the model actually gave for debugging
			raw, _ := json.Marshal(part.Parts[0])
			fmt.Println(utils.Yellow("DEBUG raw part: " + string(raw)))
			return "", 0
		}
		time.Sleep(time.Second * 7)
		return CodeGen(prompt+" give response in tool only", targetFile, allFiles, "empty rag")
	}

	// Parse function args into struct
	res, _ := json.Marshal(fc.Args)
	tkn := result.UsageMetadata.TotalTokenCount

	var suseee GenteratedCode
	if err := json.Unmarshal(res, &suseee); err != nil {
		fmt.Println(utils.Red("Failed to unmarshal function args:", err))
		return "", int32(tkn)
	}

	return suseee.Code, tkn
}
