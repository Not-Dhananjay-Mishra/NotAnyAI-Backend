package codingmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"server/models"
	"server/tools"
	"server/utils"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

type GenteratedCode struct {
	Code string `json:"code"`
}

func Itrative(data []string, prompt string, c *genai.Client, conn *websocket.Conn) map[string]string {
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
}

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
	//fmt.Println(targetFile)
	//fmt.Println(prompt)
	c := models.GeminiModel()
	ctx := context.Background()
	sysprompt := fmt.Sprintf(`
You are a frontend code assistant.

The project has React .jsx files: %v. Generate ONLY the full, valid code for "%s".

Rules:
1. Output ONLY the complete code for "%s".
2. Code must be self-contained, production-ready, responsive, modern, and engaging.
3. Use functional components with a single parent JSX element.
4. Allowed libraries: React + react-dom 18.2.0 + framer-motion 11.2.6 + Tailwind. Nothing else.
5. Do not use comments, markdown, extra text, backticks, or contractions (no dont, cant, etc).
6. JSX arrays must always use commas.
7. Apply RAG intelligently: adjust Tailwind colors, spacing, and vibe for a modern theme. RAG: %v.
8. Use image URLs only if needed, styled consistently with Tailwind and RAG.

Design & Interaction Rules:
9. Navigation handled with React state or react-dom 18.2.0 + conditional rendering (no react-router-dom).
10. Apply RAG: adjust Tailwind colors and vibes for modern/cool look. If %v is empty or irrelevant, do not force RAG.
11. Use SVG icons/illustrations (inline, optimized, accessible). No external icon libs.
12. Use animations & interactivity (hover, focus, transitions, subtle motion).
13. Create fully responsive layouts using Tailwind responsive utilities (mobile → tablet → desktop).
14. Prioritize accessibility: semantic HTML, aria-labels, keyboard navigation, sufficient color contrast.
15. Keep UI balanced: proper spacing, readable typography, clear hierarchy (header, hero, features, content, CTA, footer).
16. Avoid clutter: only include essential elements per screen.

Code Style & Structure:
17. Components must be functional React components using ES6 exports.
18. Minimal state management with useState/useEffect only.
19. Each file should contain only one component with a single responsibility.
20. Use Tailwind utilities instead of inline styles.
21. Use placeholder loading states for async content.
22. When images are needed, provide alt text and fallback.
23. SVG illustrations should be <2KB where possible and include aria-hidden or role attributes.
24. Simulate 3D with CSS/SVG unless user explicitly confirms three-js-react usage.
25. Ensure interactivity has micro-interactions (hover glow, subtle motion, button scale).
26. Include a simple footer with links, credits, or branding for completeness.

Conflict Resolution:
27. If any rules conflict, prioritize rules 4 and 5 (library and formatting constraints).
28. All outputs must strictly follow these rules without deviation.
`, allFiles, targetFile, targetFile, rag, rag)

	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{FunctionDeclarations: []*genai.FunctionDeclaration{&tools.JSXTool}},
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
		//retry without rag
		fmt.Println(utils.Red("Retrying without RAG due to error:", err))
		time.Sleep(time.Second * 70)
		res, tkn := CodeGen(prompt, targetFile, allFiles, "empty rag")
		fmt.Println(err)
		return res, tkn
	}
	if len(result.Candidates) == 0 {
		if rag == "empty rag" {
			fmt.Println(utils.Red("No candidates even after retrying without RAG"))
			return "", 0
		}
		time.Sleep(time.Second * 7)
		d, tkn := CodeGen(prompt+"give response in tool only", targetFile, allFiles, "empty rag")
		return d, tkn
	}
	if len(result.Candidates[0].Content.Parts) == 0 {
		if rag == "empty rag" {
			fmt.Println(utils.Red("No content parts even after retrying without RAG"))
			return "", 0
		}
		time.Sleep(time.Second * 7)
		d, tkn := CodeGen(prompt+"give response in tool only", targetFile, allFiles, "empty rag")
		return d, tkn
	}

	if result.Candidates[0].Content.Parts[0].FunctionCall.Args == nil {
		if rag == "empty rag" {
			fmt.Println(utils.Red("No function call args even after retrying without RAG"))
			return "", 0
		}
		res, _ := json.Marshal(result.Candidates[0].Content.Parts[0])
		fmt.Println(string(res))
		time.Sleep(time.Second * 7)
		d, tkn := CodeGen(prompt+"give response in tool only", targetFile, allFiles, "empty rag")
		return d, tkn
	}
	res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].FunctionCall.Args)
	tkn := result.UsageMetadata.TotalTokenCount
	var suseee GenteratedCode
	json.Unmarshal(res, &suseee)
	//fmt.Println(suseee.Code)
	return suseee.Code, tkn
	//return "", 0
}
