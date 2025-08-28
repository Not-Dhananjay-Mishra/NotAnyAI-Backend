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

	The frontend design consists of the following React .jsx files: %v.
	Only generate the complete, valid React .jsx code for the file "%s".

	### ABSOLUTE RULES:
	1. Output ONLY the complete code for "%s". Nothing else.
	2. The code must be self-contained, production-ready, and follow React best practices.
	3. The component must be a functional component, with all JSX wrapped in a single parent element.
	4. Use only React, react-dom, and Tailwind CSS (via className strings).
	5. Do not import or use any other external libraries (axios, classnames, react-router-dom, etc.).
	6. All JSX array elements must be separated by commas.
	7. Never output explanations, comments, markdown formatting, or extra text outside of valid code.
	8. Never output plain text or reasoning in natural language.
	9. Never use contractions like "don't" or "that's".
	10. If RAG output is provided, integrate it intelligently into the code. 
		you can change colour in tailwind classes provided by rag and put ur own color according to theame main moto is to make it modern and cool
		Do not copy-paste blindly; learn it and use accordingly into the component properly here is rag output - %v.
	11. The final output must be ONLY the valid React .jsx code, nothing else.
	12. use img url only if they needed if using then please use proper tailwind on img
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
