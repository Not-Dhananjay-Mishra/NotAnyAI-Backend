package individualtool

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"server/utils"
	"strings"
)

func GetNews(query string) []utils.NewsResult {
	newString := strings.ReplaceAll(query, " ", "%20")
	url := "https://newsapi.org/v2/everything?q=" + newString + "&sortBy=publishedAt&apiKey=" + utils.NEWS_API
	//fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching API:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil
	}
	var data utils.NewsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println("Error unmarshaling JSON:", err)
		return nil
	}

	var res []utils.NewsResult
	for _, j := range data.Articles {
		res = append(res, utils.NewsResult{
			Title:       j.Title,
			Description: j.Description,
		})
	}

	return res
}
