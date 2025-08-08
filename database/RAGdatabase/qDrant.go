package ragdatabase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

func QDrantLookup(s []float32, name string) {
	client, _ := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	sus, _ := client.Query(context.Background(), &qdrant.QueryPoints{
		CollectionName: name,
		Query:          qdrant.NewQueryDense(s),
		WithPayload:    qdrant.NewWithPayload(true)})
	res, _ := json.Marshal(sus[0].Payload)
	var response VectorSearchResponse
	json.Unmarshal(res, &response)
	fmt.Println(response.Data.Kind.StringValue)
}

func SendToVectorDB(data utils.EmbeddingResponse, collectionName string, payload string) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
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
