package individualtool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/utils"
)

const apiURL = "https://api.studio.nebius.com/v1/images/generations"

type RequestBody struct {
	Model             string      `json:"model"`
	Prompt            string      `json:"prompt"`
	ResponseFormat    string      `json:"response_format"`
	ResponseExtension string      `json:"response_extension"`
	Width             int         `json:"width"`
	Height            int         `json:"height"`
	NumSteps          int         `json:"num_inference_steps"`
	NegativePrompt    string      `json:"negative_prompt"`
	Seed              int         `json:"seed"`
	Loras             interface{} `json:"loras"`
}

type Response struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

func ImgGenHuggingFace(prompt string) string {
	reqBody := RequestBody{
		Model:             "stability-ai/sdxl",
		Prompt:            prompt,
		ResponseFormat:    "url", // ✅ change here
		ResponseExtension: "png",
		Width:             1024,
		Height:            1024,
		NumSteps:          30,
		NegativePrompt:    "",
		Seed:              -1,
		Loras:             nil,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	// Make request
	req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+utils.GEMINI_API3_IMG)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		fmt.Println("Error:", string(respBody))
		return ""
	}

	// Parse JSON response
	var apiResp Response
	_ = json.Unmarshal(respBody, &apiResp)

	if len(apiResp.Data) > 0 {
		imgURL := apiResp.Data[0].URL
		fmt.Println("✅ Image URL:", imgURL)
		return imgURL
	} else {
		fmt.Println("No image data in response")
		return ""
	}
}
