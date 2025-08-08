package codingtool

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/utils"
	"strings"
)

type SearchItem struct {
	QuestionID int    `json:"question_id"`
	Title      string `json:"title"`
	Link       string `json:"link"`
}

type SearchResponse struct {
	Items []SearchItem `json:"items"`
}

type Answer struct {
	AnswerID int    `json:"answer_id"`
	Body     string `json:"body"`
	Score    int    `json:"score"`
}

type AnswerResponse struct {
	Items []Answer `json:"items"`
}

func StackoverflowTool(query string, tag string) ([]Answer, string) {
	newquery := strings.ReplaceAll(query, " ", "+")
	searchurl := "https://api.stackexchange.com/2.3/search/advanced?order=desc&sort=relevance&q=" + newquery + "&site=stackoverflow&key=" + utils.STACKOVERFLOW_API + "&tagged=" + tag

	sresp, err := http.Get(searchurl)
	if err != nil {
		fmt.Println("Error fetching search:", err)
		return nil, ""
	}
	defer sresp.Body.Close()

	searchBytes, _ := ioutil.ReadAll(sresp.Body)

	var searchResult SearchResponse
	json.Unmarshal(searchBytes, &searchResult)

	if len(searchResult.Items) == 0 {
		fmt.Println("No matching question found.")
		return nil, ""
	}

	questionID := searchResult.Items[0].QuestionID

	answerURL := fmt.Sprintf("https://api.stackexchange.com/2.3/questions/%d/answers?order=desc&sort=activity&site=stackoverflow&filter=withbody", questionID)

	resp, err := http.Get(answerURL)
	if err != nil {
		fmt.Println("Error fetching answers:", err)
		return nil, ""
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var answerResult AnswerResponse
	json.Unmarshal(bodyBytes, &answerResult)

	if len(answerResult.Items) == 0 {
		fmt.Println("No answers found.")
		return nil, ""
	}
	return answerResult.Items, searchResult.Items[0].Title
}
