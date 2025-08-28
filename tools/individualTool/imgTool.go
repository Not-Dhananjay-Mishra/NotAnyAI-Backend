package individualtool

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"server/utils"
	"strings"
)

// Structs for parsing Pexels API response
type PexelsResponse struct {
	Photos []struct {
		Src struct {
			Original string `json:"original"`
			Large    string `json:"large"`
			Medium   string `json:"medium"`
		} `json:"src"`
		Alt string `json:"alt"`
	} `json:"photos"`
}

func ImgSearch(query string) string {

	// Prepare request
	newString := strings.ReplaceAll(query, " ", "%20")
	url := "https://api.pexels.com/v1/search?query=" + newString + "&per_page=1"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return ""
	}
	req.Header.Set("Authorization", utils.PEXELS_API_KEY)

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error fetching API:", err)
		return ""
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return ""
	}

	// Parse JSON
	var pexelsResp PexelsResponse
	if err := json.Unmarshal(body, &pexelsResp); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return ""
	}

	if len(pexelsResp.Photos) > 0 {
		// Return a nice large image
		return pexelsResp.Photos[0].Src.Large
	}

	return ""
}
