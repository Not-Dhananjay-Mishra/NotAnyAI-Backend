package codingmodel

import (
	"context"
	"encoding/json"
	"fmt"
	ragdatabase "server/database/RAGdatabase"
	"server/models"
	"server/tools"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

type GenAIResponseDB struct {
	RAG []string `json:"rag"`
	IMG []string `json:"img"`
}

const ssssprompt = `You are an AI assistant that must ONLY respond by calling the RAG tool. 
You are NOT allowed to output any freeform text, explanations, or confirmations. 
Every response must be a single tool call, nothing else. 
If the user asks for RAG Lookup queries, always generate them as exactly 7-word queries. 
Never add commentary, reasoning, or natural language text outside of the tool call.
all the RAG queries must differ from each other no same kind of queries in tool
also give search queries for img in img field keep the seach quries relevent so img can be find easily
dont give text response give only tools response
RULE - Only use tool dont give ans in text	
`

func LookupHandlerinGO(data []string) string {
	var results []string
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, i := range data {
		wg.Add(1)
		go func(payload string) {
			defer wg.Done()
			sus := LookUP(payload, "test")
			mu.Lock()
			already := false
			for _, r := range results {
				if r == sus {
					already = true
					break
				}
			}
			if !already {
				results = append(results, payload+" RAG (if required use/try dont just copy paste use acc to usecase): "+sus)
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	return strings.Join(results, " ")
}
func LookUP(pay string, coll string) string {
	c := models.GeminiModel()
	chunk := []string{pay}

	embedings := ragdatabase.DoEmbedding(c, chunk)
	if len(embedings) == 0 {
		return "No embeddings generated"
	}

	if len(embedings[0].Values) == 0 {
		return "Embedding returned with empty values"
	}

	res := ragdatabase.QDrantLookup(embedings[0].Values, coll)
	return res
}
func RAGQueryDecider(data string, conn *websocket.Conn, filename string) string {
	c := models.GeminiModel()
	ctx := context.Background()
	config := &genai.GenerateContentConfig{
		Tools: []*genai.Tool{
			{FunctionDeclarations: []*genai.FunctionDeclaration{&tools.RAG}},
		},
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: ssssprompt}}},
	}
	result, _ := c.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(data),
		config,
	)
	//log.Println(utils.Red(err))
	res, _ := json.Marshal(result.Candidates[0].Content.Parts[0].FunctionCall.Args)
	var suseee GenAIResponseDB
	json.Unmarshal(res, &suseee)
	fmt.Println(suseee.RAG)
	var images []string
	for _, i := range suseee.IMG {
		fmt.Println(i)
		//img := individualtool.ImgGenHuggingFace(i)
		//fmt.Println(utils.Magenta(img))
		//images = append(images, i+" "+img)
	}
	return LookupHandlerinGO(suseee.RAG) + " IMG " + strings.Join(images, " ")
}
