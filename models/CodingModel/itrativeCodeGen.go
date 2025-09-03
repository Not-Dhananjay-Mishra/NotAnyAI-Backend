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
	2. Code must be self-contained, production-ready, modern, with animation.
	3. Use functional components, single parent JSX, React + react-dom 18.2.0 + Tailwind only.
	4. No other libraries, no comments, no markdown, no extra text.
	5. JSX arrays need commas, no contractions.
	6. Apply RAG intelligently: adjust Tailwind colors for a modern/cool theme, not copy-paste. RAG: %v.
	7. Use img URLs only if needed, styled with proper Tailwind given in RAG.
	8. Navigation must be handled with React state or react-dom 18.2.0 and conditional rendering instead of react-router-dom.
	9. try to add animations and interactivity to make the UI more engaging.
	10. use three.js three-js-react for 3d models and animations if needed.
	11. add svg for icons and illustrations.
	12. dont use comments in the code.
	13. use images (make on ur own svg and add that) and animations to make the UI more engaging.
	14. make the home page very attractive and modern and full of content.
	`, allFiles, targetFile, targetFile, rag)

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
