package codingmodel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"server/models"
	"server/tools"
	"server/utils"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

type Process struct {
	Data    Sus
	ReplyCh chan map[string]string
	Conn    *websocket.Conn
}

var UnderProcessCode chan Process

func init() {
	UnderProcessCode = make(chan Process, 50)
}

const sysprompt = `You are a frontend code assistant. For any given prompt, 
	determine only the necessary React .jsx file names the App file must be .js App.js and 
	if required some other file in .js create that should exist to implement the described frontend design. Do not generate code, 
	explanations, or extra text—only return valid file names via the provided tool.
	The count of files should be kept to the absolute minimum required to fulfill the user's request.
	no follow up questions.
	`

type Sus struct {
	Query string `json:"query"`
}
type GenAIResponse struct {
	Filename []string `json:"file"`
}

func GetRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data Sus
	json.NewDecoder(r.Body).Decode(&data)
	replyCh := make(chan map[string]string, 1)
	UnderProcessCode <- Process{Data: data, ReplyCh: replyCh}
	result := <-replyCh
	json.NewEncoder(w).Encode(result)

}

func Processor() {
	for {
		data := <-UnderProcessCode
		if data.Data.Query != "" {
			fmt.Println(utils.Blue(data.Data.Query))
			sus := FileDecider(data.Data, data.Conn)
			data.ReplyCh <- sus
		}
	}
}

func FileDecider(data Sus, conn *websocket.Conn) map[string]string {
	c := models.GeminiModel()
	ctx := context.Background()
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{FunctionDeclarations: []*genai.FunctionDeclaration{&tools.FilenameTool}},
		},
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: sysprompt}}},
	}
	//fmt.Println("dsad")
	result, _ := c.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(data.Query),
		config,
	)
	//fmt.Println(err)
	res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].FunctionCall.Args)
	var suseee GenAIResponse
	json.Unmarshal(res, &suseee)
	fmt.Println(suseee.Filename)

	//sending statis to conn
	sendres(suseee.Filename, conn)

	//itrative generating code
	time.Sleep(time.Second * 2)
	suse := ItrativeWithoutGo(suseee.Filename, data.Query, c, conn)
	//return suse
	content := MapToContent(suse, conn)
	sus := CodingPostProcessor(content, conn, data.Query, suseee.Filename, "")
	return sus
}

func sendres(data []string, conn *websocket.Conn) {
	sus := "📂 Made files"
	for _, i := range data {
		sus = sus + " " + i
	}
	conn.WriteJSON(utils.Response{Text: sus})
}
