package utils

import (
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"google.golang.org/genai"
)

var NEWS_API string
var GOOGLE_SEARCH_API string
var GEMINI_API string
var YOUTUBE_API string

var JWT_SECRET string
var GITHUB_API string
var STACKOVERFLOW_API string
var WEATHER_API string
var GOOGLE_SEARCH_CX string

var Yellow = color.New(color.FgYellow).SprintFunc()
var Red = color.New(color.FgRed).SprintFunc()
var Green = color.New(color.FgGreen).SprintFunc()
var Cyan = color.New(color.FgCyan).SprintFunc()
var Blue = color.New(color.FgHiBlue).SprintFunc()
var Magenta = color.New(color.FgMagenta).SprintFunc()

type NewsArticle struct {
	Source      any    `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         any    `json:"url"`
}
type NewsResponse struct {
	Status       string        `json:"status"`
	TotalResults int           `json:"totalResults"`
	Articles     []NewsArticle `json:"articles"`
}

type NewsResult struct {
	Title       string
	Description string
}

type VideoDetails struct {
	Link        string
	Title       string
	ChannelName string
	View        string
	Like        string
}

type GithubRepo struct {
	FullName  string
	RepoName  string
	URL       string
	OwnerName string
}

var MemoryStore = make(map[string][]*genai.Content)

type EmbeddingResponse struct {
	Values []float32 `json:"values"`
}
type EmbeddingPayload struct {
	Values string `json:"values"`
}

var LiveConn = make(map[*websocket.Conn]bool)
var ConnUser = make(map[*websocket.Conn]string)
var UserConn = make(map[string]*websocket.Conn)

type Response struct {
	Text string `json:"text"`
}
