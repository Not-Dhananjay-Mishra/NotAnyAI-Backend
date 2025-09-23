package main

import (
	"log"
	"os"
	mongodb "server/database/MongoDB"
	codingmodel "server/models/CodingModel"
	"server/router"
	"server/utils"

	"github.com/joho/godotenv"
)

func loadEnv() {
	// Only load from .env if it exists (local dev)
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: Could not load .env file:", err)
		} else {
			log.Println("Loaded environment variables from .env")
		}
	} else {
		log.Println(".env file not found, assuming environment variables are set by the host")
	}
}

func main() {
	// Load env variables for local development
	loadEnv()
	go codingmodel.Processor()
	// Assign environment variables to utils
	utils.GEMINI_API = os.Getenv("GEMINI_API")
	utils.NEWS_API = os.Getenv("NEWS_API")
	utils.GOOGLE_SEARCH_API = os.Getenv("GOOGLE_SEARCH_API")
	utils.YOUTUBE_API = os.Getenv("YOUTUBE_API")
	utils.JWT_SECRET = os.Getenv("JWT_SECRET")
	utils.GOOGLE_SEARCH_CX = os.Getenv("GOOGLE_SEARCH_CX")
	utils.GITHUB_API = os.Getenv("GITHUB_API")
	utils.STACKOVERFLOW_API = os.Getenv("STACKOVERFLOW_API")
	utils.WEATHER_API = os.Getenv("WEATHER_API")
	utils.QDRANT_URL = os.Getenv("QDRANT_URL")
	utils.QDRANT_API = os.Getenv("QDRANT_API")
	utils.PEXELS_API_KEY = os.Getenv("PEXELS_API_KEY")
	utils.GEMINI_API3_IMG = os.Getenv("GEN_IMG")
	utils.MONGODB_USERNAME = os.Getenv("MONGODB_USERNAME")
	utils.MONGODB_PASSWORD = os.Getenv("MONGODB_PASSWORD")
	utils.MONGODB_CLUSTER = os.Getenv("MONGODB_CLUSTER")
	utils.GOOGLE_CLIENT_ID = os.Getenv("GOOGLE_CLIENT_ID")
	utils.GOOGLE_CLIENT_SECRET = os.Getenv("GOOGLE_CLIENT_SECRET")

	mongodb.DBinit()
	// Start servers
	blocker := make(chan any)
	go router.RouterHandler()

	log.Println(utils.Green("Websocket Server Started on :8000/wss/contact"))
	log.Println(utils.Green("REST API Server Started on :8000/login"))
	log.Println(utils.Green("REST API Server Started on :8000/register"))
	log.Println(utils.Green("REST API Server Started on :8000/validate"))

	<-blocker
}
