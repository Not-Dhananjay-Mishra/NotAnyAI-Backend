package ragdatabase

import (
	"context"
	"encoding/json"
	"log"
	"server/utils"

	"google.golang.org/genai"
)

func DoEmbedding(client *genai.Client, chunks []string) []utils.EmbeddingResponse {
	var ArrayEmbeddingResponse []utils.EmbeddingResponse
	var val int32 = 768
	var susint *int32 = &val
	ctx := context.Background()
	for _, chunk := range chunks {
		contents := []*genai.Content{
			genai.NewContentFromText(chunk, genai.RoleUser),
		}
		result, err := client.Models.EmbedContent(ctx,
			"gemini-embedding-001",
			contents,
			&genai.EmbedContentConfig{OutputDimensionality: susint},
		)
		if err != nil {
			log.Printf("error embedding chunk: %v\n", err)
			continue
		}

		embeddings, err := json.MarshalIndent(result.Embeddings, "", "  ")
		if err != nil {
			log.Printf("error marshalling embedding: %v\n", err)
			continue
		}
		var EmbRes []utils.EmbeddingResponse
		json.Unmarshal(embeddings, &EmbRes)
		//fmt.Println(EmbRes)
		ArrayEmbeddingResponse = append(ArrayEmbeddingResponse, EmbRes[len(EmbRes)-1])
	}
	return ArrayEmbeddingResponse
}
