package models

import (
	"fmt"
	codingtool "server/tools/codingTool"
	individualtool "server/tools/individualTool"
	websearchtool "server/tools/websearchTool"
	"server/utils"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

type allFunctionResponse struct {
	FunctionName string
	Query        string
	Response     any
}

func ToolCaller(data Agent, lastquery string, conn *websocket.Conn) []*genai.Content {
	var mu sync.Mutex
	var FunctionContent []*genai.Content
	var functiondata []allFunctionResponse
	FunctionContent = append(FunctionContent, genai.NewContentFromText(lastquery, genai.RoleUser))

	var wg sync.WaitGroup

	if data.AIGoogleSearchTool.UseTool {
		conn.WriteJSON(utils.Response{Text: "Searching on Web"})
		for _, query := range data.AIGoogleSearchTool.Query {
			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				res := websearchtool.AIGoogleSearch(q)
				//fmt.Println(utils.Yellow("AIGoogleSearchTool: "), res)
				mu.Lock()
				functiondata = append(functiondata, allFunctionResponse{FunctionName: "AIGoogleSearchTool", Query: query, Response: res})
				mu.Unlock()
			}(query)
		}
	}

	if data.DeepWikipediaSearchTool.UseTool {
		conn.WriteJSON(utils.Response{Text: "Searching on Wikipedia"})
		for _, query := range data.DeepWikipediaSearchTool.Query {

			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				res := individualtool.WikiDeepSearch(q)
				//fmt.Println(utils.Yellow("DeepWikipediaSearchTool: "), res)
				mu.Lock()
				functiondata = append(functiondata, allFunctionResponse{FunctionName: "DeepWikipediaSearchTool", Query: query, Response: res})
				mu.Unlock()
			}(query)
		}
	}

	if data.GithubSearchTool.UseTool {
		conn.WriteJSON(utils.Response{Text: "Searching on Github"})
		for _, query := range data.GithubSearchTool.Query {

			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				res := codingtool.GetRepoGithub(q)
				if len(res) > 0 {
					/*for i := range res {
						fmt.Println(utils.Yellow("GithubSearchTool: "), res[i].FullName)
					}*/
					mu.Lock()
					functiondata = append(functiondata, allFunctionResponse{FunctionName: "GithubSearchTool", Query: query, Response: res})
					mu.Unlock()
				}
			}(query)
		}
	}

	if data.GoogleSearchTool.UseTool {
		conn.WriteJSON(utils.Response{Text: "Searching on Web"})
		for _, query := range data.GoogleSearchTool.Query {
			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				link, _ := websearchtool.RESTGoogleSearch(q)
				if len(link) > 0 {
					/*for i := range link {
						fmt.Println(utils.Yellow("RESTGoogleSearch: "), link[i])
					}*/
					mu.Lock()
					functiondata = append(functiondata, allFunctionResponse{FunctionName: "GoogleSearchTool", Query: query, Response: map[string]any{"link": link, "content": data}})
					mu.Unlock()
				}
			}(query)
		}
	}

	if data.NewsSearchTool.UseTool {
		conn.WriteJSON(utils.Response{Text: "Searching for News"})
		for _, query := range data.NewsSearchTool.Query {

			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				res := individualtool.GetNews(q)
				if len(res) > 0 {
					/*for i := range res {
						fmt.Println(utils.Yellow("NewsSearchTool: "), res[i].Description)

					}*/
					mu.Lock()
					functiondata = append(functiondata, allFunctionResponse{FunctionName: "NewsSearchTool", Query: query, Response: res})
					mu.Unlock()
				}
			}(query)
		}
	}

	if data.RedditSearchTool.UseTool {
		conn.WriteJSON(utils.Response{Text: "Searching on Reddit"})
		for _, query := range data.RedditSearchTool.Query {
			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				fmt.Println(utils.Red("Reddit Not made yet for query:"), q)

			}(query)
		}
	}

	if data.StackoverflowSearchTool.UseTool {
		conn.WriteJSON(utils.Response{Text: "Searching on Stackoverflow"})
		for _, query := range data.StackoverflowSearchTool.Query {
			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				fmt.Println(utils.Red("Stackoverflow Not made yet for query:"), q)
			}(query)
		}
	}

	if data.WeatherTool.UseTool {
		conn.WriteJSON(utils.Response{Text: "Fetching current weather"})
		for _, query := range data.WeatherTool.Query {

			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				name, temp, description := individualtool.WeatherTool(q)
				//fmt.Println(utils.Yellow("WeatherTool: "), name, temp, description)
				mu.Lock()
				functiondata = append(functiondata, allFunctionResponse{FunctionName: "WeatherTool", Query: query, Response: map[string]any{"place": name, "temp": temp, "description": description}})
				mu.Unlock()
			}(query)
		}
	}

	if data.WikipediaSearchTool.UseTool {
		conn.WriteJSON(utils.Response{Text: "Searching on Wikipedia"})
		for _, query := range data.WikipediaSearchTool.Query {

			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				res := individualtool.WikiSummarySearch(q)
				//fmt.Println(utils.Yellow("WikipediaSearchTool: "), res)
				mu.Lock()
				functiondata = append(functiondata, allFunctionResponse{FunctionName: "WikipediaSearchTool", Query: query, Response: res})
				mu.Unlock()
			}(query)
		}
	}

	if data.YoutubePlaylistTool.UseTool {
		conn.WriteJSON(utils.Response{Text: "Searching on Youtube"})
		for _, query := range data.YoutubePlaylistTool.Query {

			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				link, title, descriptions := individualtool.YoutubeToolPlaylist(q)
				if len(link) > 0 && len(title) > 0 && len(descriptions) > 0 {
					/*for i := range link {
						fmt.Println(utils.Yellow("YoutubePlaylistTool: "), link[i], title[i], descriptions[i])
					}*/
					mu.Lock()
					functiondata = append(functiondata, allFunctionResponse{FunctionName: "YoutubePlaylistTool", Query: query, Response: map[string]any{"link": link, "title": title, "description": descriptions}})
					mu.Unlock()
				}
			}(query)
		}
	}

	if data.YoutubeVideoTool.UseTool {
		conn.WriteJSON(utils.Response{Text: "Searching on Youtube"})
		for _, query := range data.YoutubeVideoTool.Query {

			wg.Add(1)
			go func(q string) {
				defer wg.Done()
				res := individualtool.YoutubeToolVideo(q)
				if len(res) > 0 {
					/*for i := range res {
						fmt.Println(utils.Yellow("YoutubeVideoTool: "), res[i].Title, res[i].Link)
					}*/
					mu.Lock()
					functiondata = append(functiondata, allFunctionResponse{FunctionName: "YoutubeVideoTool", Query: query, Response: res})
					mu.Unlock()
				}
			}(query)
		}
	}

	wg.Wait()
	//fmt.Println(utils.Green(functiondata))
	for _, out := range functiondata {
		FunctionContent = append(FunctionContent, genai.NewContentFromFunctionCall(out.FunctionName, map[string]any{
			"query": out.Query,
		}, genai.RoleModel))

		FunctionContent = append(FunctionContent, genai.NewContentFromFunctionResponse(out.FunctionName, map[string]any{
			out.Query: out.Response,
		}, genai.RoleModel))
	}
	return FunctionContent
}
