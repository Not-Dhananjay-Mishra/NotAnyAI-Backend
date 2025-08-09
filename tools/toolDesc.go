package tools

import "google.golang.org/genai"

var WikipediaSearchToolDescription = genai.FunctionDeclaration{
	Name:        "wikipedia_search",
	Description: "Searches Wikipedia for information based on the given query. summary",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"query": {
				Type:        genai.TypeArray,
				Description: "query for searching",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"query"},
	},
}
var NewsSearchToolDescription = genai.FunctionDeclaration{
	Name:        "news_search",
	Description: "Searches News handles for information based on the given query.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"query": {
				Type:        genai.TypeArray,
				Description: "query for searching",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"query"},
	},
}
var DeepWikipediaSearchToolDescription = genai.FunctionDeclaration{
	Name:        "wikipedia_deep_search",
	Description: "Searches Deep Wikipedia for information based on the given query. in depth",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"query": {
				Type:        genai.TypeArray,
				Description: "query for searching",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"query"},
	},
}
var RedditSearchToolDescription = genai.FunctionDeclaration{
	Name:        "reddit_search",
	Description: "Searches on reddit for infomation",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"subreddit": {
				Type:        genai.TypeArray,
				Description: "subreddit name for searching",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
			"query": {
				Type:        genai.TypeArray,
				Description: "query for searching",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"subreddit", "query"},
	},
}
var WeatherToolDescription = genai.FunctionDeclaration{
	Name:        "weather_search",
	Description: "gives weather details of any given place",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"place": {
				Type:        genai.TypeArray,
				Description: "place name whose weather is required",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"place"},
	},
}
var YoutubeVideoToolDescription = genai.FunctionDeclaration{
	Name:        "youtube_video",
	Description: "gives youtube video for given search",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"query": {
				Type:        genai.TypeArray,
				Description: "query to search on youtube",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"query"},
	},
}
var YoutubePlaylistToolDescription = genai.FunctionDeclaration{
	Name:        "youtube_playlist",
	Description: "gives youtube playlist for given search",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"query": {
				Type:        genai.TypeArray,
				Description: "query to search on youtube",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"query"},
	},
}
var GoogleSearchToolDescription = genai.FunctionDeclaration{
	Name:        "google_search",
	Description: "search the web and retrive url if required gives html code too",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"query": {
				Type:        genai.TypeArray,
				Description: "query to search on google",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
			"wantHTML": {
				Type:        genai.TypeBoolean,
				Description: "want HTML for query",
			},
		},
		Required: []string{"query", "wantHTML"},
	},
}
var AIGoogleSearchToolDescription = genai.FunctionDeclaration{
	Name:        "google_search_ai",
	Description: "search the web and gives summary using genai",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"query": {
				Type:        genai.TypeArray,
				Description: "query to search on google",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"query"},
	},
}
var StackoverflowSearchToolDescription = genai.FunctionDeclaration{
	Name:        "stackoverflow_search",
	Description: "search for any coding problem/bug on stackoverflow",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"query": {
				Type:        genai.TypeArray,
				Description: "query to search on stackoverflow",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
			"tag": {
				Type:        genai.TypeArray,
				Description: "tag for the queries",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"query", "tag"},
	},
}
var GithubSearchToolDescription = genai.FunctionDeclaration{
	Name:        "github_search",
	Description: "Search for repositories, code, or documentation on GitHub. This is useful for answering questions about how to use libraries or technologies, especially in a specific programming language.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"query": {
				Type:        genai.TypeArray,
				Description: "query to search github",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
			"wantReadme": {
				Type:        genai.TypeBoolean,
				Description: "want readme file too",
			},
		},
		Required: []string{"query", "wantReadme"},
	},
}

var ToolDeciderAgent = genai.FunctionDeclaration{
	Name:        "agent",
	Description: "Decides which tools to use based on the query u can set true to all relevent tools",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"wikipediatool": {
				Type:        genai.TypeObject,
				Description: "Configuration for Wikipedia tools.",
				Properties: map[string]*genai.Schema{
					"usetool": {
						Type:        genai.TypeBoolean,
						Description: "Set to true if the Wikipedia tool should be used.",
					},
					"query": {
						Type:        genai.TypeArray,
						Description: "An array of independent search queries extracted from user intent give good context for queries",
						Items:       &genai.Schema{Type: genai.TypeString},
					},
				},
			},
			"deepwikipediatool": {
				Type:        genai.TypeObject,
				Description: "Configuration for Deep Wikipedia tools.",
				Properties: map[string]*genai.Schema{
					"usetool": {
						Type:        genai.TypeBoolean,
						Description: "Set to true if the Deep Wikipedia tool should be used.",
					},
					"query": {
						Type:        genai.TypeArray,
						Description: "An array of independent search queries extracted from user intent",
						Items:       &genai.Schema{Type: genai.TypeString},
					},
				},
			},
			"news_searchtool": {
				Type:        genai.TypeObject,
				Description: "Configuration for News Search tools.",
				Properties: map[string]*genai.Schema{
					"usetool": {
						Type:        genai.TypeBoolean,
						Description: "Set to true if recent news articles should be searched.",
					},
					"query": {
						Type:        genai.TypeArray,
						Description: "An array of independent search queries extracted from user intent",
						Items:       &genai.Schema{Type: genai.TypeString},
					},
				},
			},
			"github_searchtool": {
				Type:        genai.TypeObject,
				Description: "Configuration for Github Search tools.",
				Properties: map[string]*genai.Schema{
					"usetool": {
						Type:        genai.TypeBoolean,
						Description: "Set to true if GitHub should be searched for repositories or code.",
					},
					"query": {
						Type:        genai.TypeArray,
						Description: "An array of independent search queries extracted from user intent",
						Items:       &genai.Schema{Type: genai.TypeString},
					},
				},
			},
			"stackoverflow_searchtool": {
				Type:        genai.TypeObject,
				Description: "Configuration for Stackoverflow Search tools.",
				Properties: map[string]*genai.Schema{
					"usetool": {
						Type:        genai.TypeBoolean,
						Description: "Set to true if Stack Overflow should be searched for programming questions.",
					},
					"query": {
						Type:        genai.TypeArray,
						Description: "An array of independent search queries extracted from user intent",
						Items:       &genai.Schema{Type: genai.TypeString},
					},
				},
			},
			"google_search_aitool": {
				Type:        genai.TypeObject,
				Description: "search the web and give infomation of pages. better",
				Properties: map[string]*genai.Schema{
					"usetool": {
						Type:        genai.TypeBoolean,
						Description: "Set to true if AI-powered Google search tool should be used.",
					},
					"query": {
						Type:        genai.TypeArray,
						Description: "An array of independent search queries extracted from user intent",
						Items:       &genai.Schema{Type: genai.TypeString},
					},
				},
			},
			"google_searchtool": {
				Type:        genai.TypeObject,
				Description: "search the web and give infomation of pages.",
				Properties: map[string]*genai.Schema{
					"usetool": {
						Type:        genai.TypeBoolean,
						Description: "Set to true if standard Google search should be used.",
					},
					"query": {
						Type:        genai.TypeArray,
						Description: "An array of independent search queries extracted from user intent",
						Items:       &genai.Schema{Type: genai.TypeString},
					},
				},
			},
			"youtube_playlisttool": {
				Type:        genai.TypeObject,
				Description: "Configuration for youtube playlist search tools.",
				Properties: map[string]*genai.Schema{
					"usetool": {
						Type:        genai.TypeBoolean,
						Description: "Set to true if a YouTube playlist should be searched.",
					},
					"query": {
						Type:        genai.TypeArray,
						Description: "An array of independent search queries extracted from user intent",
						Items:       &genai.Schema{Type: genai.TypeString},
					},
				},
			},
			"youtube_videotool": {
				Type:        genai.TypeObject,
				Description: "Configuration for youtube video search tools.",
				Properties: map[string]*genai.Schema{
					"usetool": {
						Type:        genai.TypeBoolean,
						Description: "Set to true if a YouTube video should be searched.",
					},
					"query": {
						Type:        genai.TypeArray,
						Description: "An array of independent search queries extracted from user intent",
						Items:       &genai.Schema{Type: genai.TypeString},
					},
				},
			},
			"weathertool": {
				Type:        genai.TypeObject,
				Description: "Configuration for live weather tools.",
				Properties: map[string]*genai.Schema{
					"usetool": {
						Type:        genai.TypeBoolean,
						Description: "Set to true if current weather data should be fetched.",
					},
					"query": {
						Type:        genai.TypeArray,
						Description: "An array of independent place name only extracted from user intent",
						Items:       &genai.Schema{Type: genai.TypeString},
					},
				},
			},
			"reddit_searchtool": {
				Type:        genai.TypeObject,
				Description: "Configuration for reddit search tools.",
				Properties: map[string]*genai.Schema{
					"usetool": {
						Type:        genai.TypeBoolean,
						Description: "Set to true if Reddit should be searched for discussions or opinions.",
					},
					"query": {
						Type:        genai.TypeArray,
						Description: "An array of independent search queries extracted from user intent",
						Items:       &genai.Schema{Type: genai.TypeString},
					},
				},
			},
		},
		//Required: []string{"reddit_searchtool"},
	},
}
