package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/tools"
	"server/utils"

	"google.golang.org/genai"
)

type FunctionArgs struct {
	Query   []string `json:"query"`
	UseTool bool     `json:"usetool"`
}

type Agent struct {
	WikipediaSearchTool     FunctionArgs `json:"wikipediatool"`
	NewsSearchTool          FunctionArgs `json:"news_searchtool"`
	DeepWikipediaSearchTool FunctionArgs `json:"deepwikipediatool"`
	RedditSearchTool        FunctionArgs `json:"reddit_searchtool"`
	WeatherTool             FunctionArgs `json:"weathertool"`
	YoutubeVideoTool        FunctionArgs `json:"youtube_videotool"`
	YoutubePlaylistTool     FunctionArgs `json:"youtube_playlisttool"`
	GoogleSearchTool        FunctionArgs `json:"google_searchtool"`
	AIGoogleSearchTool      FunctionArgs `json:"google_search_aitool"`
	StackoverflowSearchTool FunctionArgs `json:"stackoverflow_searchtool"`
	GithubSearchTool        FunctionArgs `json:"github_searchtool"`
}

func ModelWithTools(c *genai.Client, prompt []*genai.Content, username string) string {
	ctx := context.Background()
	if len(prompt) == 0 || len(prompt[len(prompt)-1].Parts) == 0 {
		log.Println("Prompt or parts is empty")
		return "Invalid prompt"
	}
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
	fmt.Println(utils.Magenta("Prompt: "), prompt[len(prompt)-1].Parts[0].Text)
	result, err := c.Models.GenerateContent(
		ctx,
		"gemini-2.5-pro",
		prompt,
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
	if err != nil {
		log.Fatal(err)
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		res, _ := json.Marshal(result)
		log.Println(string(res))
		for i := range prompt {
			log.Println(prompt[i].Parts)
		}
		log.Println("No candidates or parts returned from model")
		return "Sorry, I couldn't find any information for that."
	}

	part := result.Candidates[0].Content.Parts[0]

	if part.Text != "" {
		utils.MemoryStore[username] = append(utils.MemoryStore[username], genai.NewContentFromText(part.Text, genai.RoleModel))
		ans := ""
		for i := range result.Candidates[0].Content.Parts {
			//res, _ := json.Marshal(result.Candidates[0].Content.Parts[i].Text)
			ans += result.Candidates[0].Content.Parts[i].Text
		}
		fmt.Println(utils.Yellow("AI : "), ans)
		return ans
	} else if part.FunctionCall.Name != "" {
		res, _ := json.Marshal(part.FunctionCall.Args)
		var data Agent
		json.Unmarshal(res, &data)
		fmt.Println(utils.Cyan(string(res)))
		content := ToolCaller(data, prompt[len(prompt)-1].Parts[0].Text)
		sus := PostProcessing(c, username, content, prompt[len(prompt)-1].Parts[0].Text, prompt)
		return sus
	} else {
		log.Println("Error processing request: No text or function call")
		return "Error processing request"
	}
}
