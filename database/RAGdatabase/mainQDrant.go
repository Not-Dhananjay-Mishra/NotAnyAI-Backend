package ragdatabase

import (
	"context"
	"log"
	"server/utils"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

func RAGinit(username string, data []utils.EmbeddingResponse, chunks []string) string {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   utils.QDRANT_URL,
		Port:   6334,
		APIKey: utils.QDRANT_API,
	})
	if err != nil {
		log.Println(utils.Red("Error Connecting to Qdrant:"), err)
		return ""
	}
	ctx := context.Background()

	collectionName := username + uuid.New().String()
	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(768),
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		log.Println("Collection creation error:", err)
		return ""
	}
	for i := range data {
		SendToVectorDB(data[i], collectionName, chunks[i], chunks[i])
	}
	return collectionName
}
