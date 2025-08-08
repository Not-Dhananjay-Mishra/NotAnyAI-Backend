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
		Host: "localhost",
		Port: 6334,
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
		SendToVectorDB(data[i], collectionName, chunks[i])
	}
	return collectionName
}
