package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/tools"
	"server/utils"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

func ImageModel(client *genai.Client, path string, prompt string, imgFormat string, conn *websocket.Conn, username string, imgbyte []byte) string {
	conn.WriteJSON(utils.Response{Text: "Thinking..."})
	/*bytes, err := os.ReadFile(path)
	if err != nil {
		bytes = imgbyte
	}*/
	bytes := imgbyte
	ctx := context.Background()
	imgtype := "image/" + imgFormat
	sus := `
						You are a helpful AI assistant.
						Your primary goal is to give the user the most complete and accurate answer possible.

						Rules:
						1. If you can answer directly with your own knowledge, do so without calling tools.
						2. Only call tools when the answer requires external, real-time, or highly specific information.
						3. When giving code, return the **entire working code block** in the requested language without extra setup steps unless explicitly asked.
						4. Maintain proper formatting and preserve line breaks in code.
						5. You may ask clarifying questions if the prompt is ambiguous.
						6. Always ensure the final output is ready for the user to copy and use immediately.
						7. if using only llm give full response
						Follow these rules strictly.
			`

	parts := []*genai.Part{
		genai.NewPartFromBytes(bytes, imgtype),
		genai.NewPartFromText(prompt),
	}

	contents := []*genai.Content{
		genai.NewContentFromParts(parts, genai.RoleUser),
	}

	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		contents,
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				{
					FunctionDeclarations: []*genai.FunctionDeclaration{
						&tools.ToolDeciderAgent},
				},
			},
			SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: sus}}},
		},
	)
	finalans := ""

	if result.Candidates[0].Content.Parts[0].Text != "" {
		//fmt.Println(result.Candidates[0].Content.Parts[0].Text)
		finalans = result.Candidates[0].Content.Parts[0].Text
	} else if len(result.Candidates[0].Content.Parts[0].FunctionCall.Args) != 0 {
		res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].FunctionCall.Args)
		var sus Agent
		json.Unmarshal(res, &sus)
		log.Println(utils.Blue(string(res)))
		content := ImgToolCaller(sus, parts, conn)
		//fmt.Println(content)
		//PrintMemobyContent(content)
		finalans += ImgStreamPostProcessing(client, username, content, prompt, contents, conn)

	} else {
		res, _ := json.Marshal(result)
		fmt.Println(utils.Red(res))
		return "Error"
	}

	//utils.MemoryStore[username] = append(utils.MemoryStore[username], genai.NewContentFromText(finalans, genai.RoleModel))
	return finalans
}
