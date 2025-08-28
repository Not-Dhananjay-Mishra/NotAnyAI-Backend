package ragdatabase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"server/utils"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

type VectorSearchResponse struct {
	Data struct {
		Kind struct {
			StringValue string `json:"StringValue"`
		} `json:"kind"`
	} `json:"data"`
}

func QDrantLookup(s []float32, name string) string {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   utils.QDRANT_URL,
		Port:   6334,
		APIKey: utils.QDRANT_API,
		UseTLS: true,
	})
	if err != nil {
		return ""
	}

	points, err := client.Scroll(context.Background(), &qdrant.ScrollPoints{
		CollectionName: name,
		WithPayload:    qdrant.NewWithPayload(true),
		Limit:          utils.ToPtr(uint32(10)), // fetch 10 candidates
	})
	if err != nil || len(points) == 0 {
		return ""
	}

	idx := rand.Intn(len(points)) // pick random from fetched batch

	res, _ := json.Marshal(points[idx].Payload)
	var response VectorSearchResponse
	json.Unmarshal(res, &response)
	return response.Data.Kind.StringValue
}

func SendToVectorDB(data utils.EmbeddingResponse, collectionName string, payload string, emb string) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   utils.QDRANT_URL,
		Port:   6334,
		APIKey: utils.QDRANT_API,
		UseTLS: true,
	})
	if err != nil {
		log.Println(utils.Red("Error Connecting to Qdrant:"), err)
		return
	}

	ctx := context.Background()

	// ONCEEEEEEEEEE..
	/*err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(len(data.Values)),
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		log.Println("Collection creation error:", err)
		return
	}*/

	// Insert vectors
	vec := data.Values
	op, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDUUID(uuid.New().String()),
				Vectors: qdrant.NewVectorsDense(vec),
				Payload: qdrant.NewValueMap(map[string]any{"data": payload}),
			},
		},
	})
	if err != nil {
		log.Println("Upsert error at index", ":", err)
	} else {
		fmt.Println("Upsert response for index", ":", op)
	}

}
