package websearchtool

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"google.golang.org/genai"
)

func Chunker(data string) []string {
	const minChunkSize = 500
	length := len(data)

	if length <= minChunkSize {
		return []string{data}
	}

	chunksCount := length / minChunkSize
	if length%minChunkSize != 0 {
		chunksCount++
	}

	chunkSize := length / chunksCount
	remainder := length % chunksCount

	result := make([]string, 0, chunksCount)
	start := 0

	for i := 0; i < chunksCount; i++ {
		end := start + chunkSize
		if remainder > 0 {
			end++
			remainder--
		}
		if end > length {
			end = length
		}
		result = append(result, data[start:end])
		start = end
	}
	return result
}

func Scrapper(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch %s: status code %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func AIScrapper(c *genai.Client, url1 string) string {
	ctx := context.Background()

	sus := fmt.Sprintf("Retrive all the content of the provided URL. %v ", url1)
	//fmt.Println(utils.Blue(sus))

	prompt := genai.NewContentFromParts([]*genai.Part{
		{Text: sus},
	}, genai.RoleModel)
	result, err := c.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		[]*genai.Content{prompt},
		&genai.GenerateContentConfig{
			Tools: []*genai.Tool{
				{
					//URLContext:   &genai.URLContext{},
					GoogleSearch: &genai.GoogleSearch{},
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Print the model's response.
	res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].Text)
	//fmt.Println(utils.Blue(string(res)))
	return (string(res))
}
