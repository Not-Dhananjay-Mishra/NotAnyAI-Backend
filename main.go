package main

import (
	"log"
	"os"
	"server/router"
	"server/utils"

	"github.com/joho/godotenv"
)

func main() {
	/*username := "techdm"
	client := models.GeminiModel()
	prompt := "tell me about latest india news"
	models.AddToMemoryUSER(username, prompt)
	models.ModelWithTools(client, utils.MemoryStore[username], username)

	prompt = "what can be its major impact"
	models.AddToMemoryUSER(username, prompt)
	models.ModelWithTools(client, utils.MemoryStore[username], username)*/
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//Loading env
	utils.GEMINI_API = os.Getenv("GEMINI_API")
	utils.NEWS_API = os.Getenv("NEWS_API")
	utils.GOOGLE_SEARCH_API = os.Getenv("GOOGLE_SEARCH_API")
	utils.YOUTUBE_API = os.Getenv("YOUTUBE_API")
	utils.JWT_SECRET = os.Getenv("JWT_SECRET")
	utils.GOOGLE_SEARCH_CX = os.Getenv("GOOGLE_SEARCH_CX")
	utils.GITHUB_API = os.Getenv("GITHUB_API")
	utils.STACKOVERFLOW_API = os.Getenv("STACKOVERFLOW_API")
	utils.WEATHER_API = os.Getenv("WEATHER_API")

	blocker := make(chan any)

	go router.RouterHandler()
	log.Println(utils.Green("Websocket Server Started on localhost:8000/wss/contact"))
	log.Println(utils.Green("REST API Server Started on localhost:8000/login"))
	log.Println(utils.Green("REST API Server Started on localhost:8000/register"))
	log.Println(utils.Green("REST API Server Started on localhost:8000/validate"))

	<-blocker
}
