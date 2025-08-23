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
			//conn.WriteJSON(utils.Response{Text: "⚙️ Started generation of " + file})
			conn.WriteJSON(map[string]string{"codegenstart": "⚙️ Started generation of " + file})
			code, tkn := CodeGen(prompt, fname, data)
			//conn.WriteJSON(utils.Response{Text: "✅ Completed generation of " + file})
			conn.WriteJSON(map[string]string{"codegencomplete": "✅ Completed generation of " + file})
			mu.Lock()
			results[fname] = code
			totaltkn += int(tkn)
			mu.Unlock()
		}(file)
		time.Sleep(time.Second * 8)
	}
	wg.Wait()
	fmt.Println(utils.Red("Code Gen Token - ", totaltkn))
	return results
}

func ItrativeWithoutGo(data []string, prompt string, c *genai.Client, conn *websocket.Conn) map[string]string {
	// without go rountines
	results := make(map[string]string)
	totaltkn := 0

	for _, file := range data {
		conn.WriteJSON(utils.Response{Text: "⚙️ Started generation of " + file})
		code, tkn := CodeGen(prompt, file, data)
		conn.WriteJSON(utils.Response{Text: "✅ Completed generation of " + file})
		results[file] = code
		totaltkn += int(tkn)
		time.Sleep(time.Second * 2)
	}
	fmt.Println(utils.Red("Code Gen Token - ", totaltkn))
	return results
}

func CodeGen(prompt string, targetFile string, allFiles []string) (string, int32) {
	//fmt.Println(targetFile)
	//fmt.Println(prompt)
	c := models.GeminiModel()
	ctx := context.Background()
	sysprompt := fmt.Sprintf(`
	You are a frontend code assistant.

	The frontend design consists of the following React .jsx files: %v.
	Only generate the complete, valid React .jsx code for the file "%s".

	### Rules:
	1. The code must be self-contained, production-ready, and follow React best practices.
	2. The component must be a functional component and the returned JSX must be wrapped in a single parent element.
	3. Use only React, react-dom and Tailwind CSS (via className strings).
	4. Do not use any external libraries (e.g., axios, classnames, react-router-dom (alternative of this can be used - react-dom), etc.).
	5. All JSX array elements must be separated by commas.
	6. Do not output explanations, comments, extra text, backslashes (\), or forward slashes (/) unless absolutely required inside valid JSX or string literals.
	7. Do not use contractions like "don't" or "that's etc that has : ' in it".
	8. Only return the code content for "%s".
	`, allFiles, targetFile, targetFile)

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
		fmt.Println(utils.Red("TOKEN KATAM HOGAYE!!"))

		fmt.Println(err)
		return "", 0
	}
	if result.Candidates[0].Content.Parts[0].FunctionCall.Args == nil {
		res, _ := json.Marshal(result.Candidates[0].Content.Parts[0])
		fmt.Println(string(res))
		return "", result.UsageMetadata.TotalTokenCount
	}
	res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].FunctionCall.Args)
	tkn := result.UsageMetadata.TotalTokenCount
	var suseee GenteratedCode
	json.Unmarshal(res, &suseee)
	//fmt.Println(suseee.Code)
	return suseee.Code, tkn
	//return "", 0
}
