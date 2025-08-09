package websearchtool

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"server/utils"
	"sync"

	"google.golang.org/genai"
)

func RESTGoogleSearch(query string) ([]string, []string) {
	var wg sync.WaitGroup
	baseURL := "https://www.googleapis.com/customsearch/v1"
	params := url.Values{}
	params.Add("key", utils.GOOGLE_SEARCH_API)
	params.Add("cx", utils.GOOGLE_SEARCH_CX)
	params.Add("q", query)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()

	var result struct {
		Items []struct {
			Title string `json:"title"`
			Link  string `json:"link"`
		} `json:"items"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, nil
	}

	var links []string
	var data []string
	ctx := context.Background()
	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  utils.GEMINI_API,
		Backend: genai.BackendGeminiAPI,
	})
	//log.Println("url done")
	if len(result.Items) > 2 {
		for i := 0; i < 2; i++ {
			links = append(links, result.Items[i].Link)
			//scrappeddata, _ := Scrapper(result.Items[i].Link)
			//data = append(data, scrappeddata)
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				data = append(data, AIScrapper(client, url))
			}(result.Items[i].Link)

		}
		wg.Wait()

	}

	return links, data
}

func AIGoogleSearch(query string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  utils.GEMINI_API,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}
	r, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(query),
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				{
					GoogleSearch: &genai.GoogleSearch{},
				},
			},
			SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: "You can answer questions by using only tools"}}},
		},
	)
	if err != nil {
		log.Println(err)
	}
	tkn, _ := json.Marshal(r.UsageMetadata.TotalTokenCount)
	res, _ := json.Marshal(r.Candidates[0].Content.Parts[0])
	fmt.Println("Total Token used: ", string(tkn))
	return string(res)

}
