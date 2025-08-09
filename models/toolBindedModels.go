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
	fmt.Println(utils.Magenta("Prompt: "), prompt[len(prompt)-1].Parts[0].Text)
	result, err := c.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		prompt,
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				{
					FunctionDeclarations: []*genai.FunctionDeclaration{
						&tools.ToolDeciderAgent},
				},
			},
			SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: "You are a helpful assistant. You can answer questions normally or use tools if required. u can ask for follow up"}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	tkn, _ := json.Marshal(result.UsageMetadata.TotalTokenCount)
	fmt.Println("Total Token used: ", string(tkn))
	if result.Candidates[0].Content.Parts[0].Text != "" {
		res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].Text)
		utils.MemoryStore[username] = append(utils.MemoryStore[username], genai.NewContentFromText(string(res), genai.RoleModel))
		fmt.Println(utils.Yellow("AI : "), string(res))
		return string(res)
	} else if result.Candidates[0].Content.Parts[0].FunctionCall.Name != "" {
		res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].FunctionCall.Args)
		var data Agent
		json.Unmarshal(res, &data)
		fmt.Println(utils.Cyan(string(res)))
		content := ToolCaller(data, prompt[len(prompt)-1].Parts[0].Text)
		sus := PostProcessing(c, username, content, prompt[len(prompt)-1].Parts[0].Text, prompt)
		return sus
	} else {
		log.Println("Error Proscessing request")
		return "Error Proscessing request"
	}
}
