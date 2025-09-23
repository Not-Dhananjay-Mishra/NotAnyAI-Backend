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
	ReplyCh chan PostCodeResponse
	Conn    *websocket.Conn
}

var UnderProcessCode chan Process

func init() {
	UnderProcessCode = make(chan Process, 50)
}

const sysprompt = `You are a frontend and backend code assistant for generating Next.js projects.
For any given prompt, determine only the necessary files required to implement the described website.
- Frontend pages must be returned in "frontendFiles" as .js files (e.g., pages/index.js, pages/about.js).
- Backend API endpoints must be returned in "backendFiles" as .js files under pages/api/ (e.g., pages/api/hello.js).
- Do not generate code, explanations, or extra text—only return valid file names via the provided tool.
- Keep the count of files to the absolute minimum required to fulfill the user's request.
- dont create nested file execpt pages and api only /pages/<filename>.js and /pages/api/<filename>.js nothing other than that
- No follow-up questions.`

type Sus struct {
	Query string `json:"query"`
}
type GenAIResponse struct {
	FrontendFile []string `json:"frontendFiles"`
	BackendFile  []string `json:"backendFiles"`
}

func GetRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data Sus
	json.NewDecoder(r.Body).Decode(&data)
	replyCh := make(chan PostCodeResponse, 1)
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

func FileDecider(data Sus, conn *websocket.Conn) PostCodeResponse {
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
	fmt.Println(utils.Yellow(data.Query))
	//fmt.Println(err)
	res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].FunctionCall.Args)
	var suseee GenAIResponse
	json.Unmarshal(res, &suseee)
	fmt.Println(suseee)

	//sending statis to conn
	if len(suseee.FrontendFile) > 0 {
		sendres(suseee.FrontendFile, conn)
	}
	if len(suseee.BackendFile) > 0 {
		sendres(suseee.BackendFile, conn)
	}

	//itrative generating code
	time.Sleep(time.Second * 2)
	var AllFiles []string
	AllFiles = append(AllFiles, suseee.BackendFile...)
	AllFiles = append(AllFiles, suseee.FrontendFile...)

	suse := ItrativeWithoutGo(AllFiles, data.Query, c, conn)
	//return suse
	content := MapToContent(suse, conn)
	sus := CodingPostProcessor(content, conn, data.Query, AllFiles, "")
	return sus
}

func sendres(data []string, conn *websocket.Conn) {
	sus := "📂 Made files"
	for _, i := range data {
		sus = sus + " " + i
	}
	conn.WriteJSON(utils.Response{Text: sus})
}
