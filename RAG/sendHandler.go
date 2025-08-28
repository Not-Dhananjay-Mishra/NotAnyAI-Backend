package rag

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	ragdatabase "server/database/RAGdatabase"
	"server/models"
	"server/utils"

	"github.com/qdrant/go-client/qdrant"
)

type RequestPayload struct {
	EmbeddingPayload string `json:"EmbeddingPayload"`
	Payload          string `json:"payload"`
	Collection       string `json:"collection"`
}

type RequestLookup struct {
	Payload    string `json:"payload"`
	Collection string `json:"collection"`
}

func RAGSend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data RequestPayload
	var chunk []string
	json.NewDecoder(r.Body).Decode(&data)
	//log.Println("Request EmbeddingPayload: ", data.EmbeddingPayload)
	chunk = append(chunk, data.EmbeddingPayload)
	c := models.GeminiModel()

	embedings := ragdatabase.DoEmbedding(c, chunk)

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

	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: data.Collection,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     uint64(768),
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		log.Println("Already Exist")
	}
	ragdatabase.SendToVectorDB(embedings[0], data.Collection, data.Payload, data.EmbeddingPayload)
	log.Println(utils.Green("Done!"))
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "done"})
}

func RAGLookup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data RequestLookup
	var chunk []string
	json.NewDecoder(r.Body).Decode(&data)
	chunk = append(chunk, data.Payload)
	c := models.GeminiModel()

	embedings := ragdatabase.DoEmbedding(c, chunk)

	res := ragdatabase.QDrantLookup(embedings[0].Values, data.Collection)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"find": res})
}
