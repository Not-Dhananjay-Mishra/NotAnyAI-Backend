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
			"usesitecraft": {
				Type:        genai.TypeObject,
				Description: "it makes website/web app using AI. only use this tool when user ask for website",
				Properties: map[string]*genai.Schema{
					"usetool": {
						Type:        genai.TypeBoolean,
						Description: "Swt to True if want to make website.",
					},
					"query": {
						Type:        genai.TypeArray,
						Description: "queries for makeing website",
						Items:       &genai.Schema{Type: genai.TypeString},
					},
				},
			},
		},
		//Required: []string{"reddit_searchtool"},
	},
}

// temp new try
var NewToolDeciderAgent = genai.FunctionDeclaration{
	Name:        "agent",
	Description: "Select all relevant tools for the user query. Each tool has a 'usetool' boolean and a 'query' array of independent, context-rich search queries extracted from user intent.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"wikipediatool": {
				Type:        genai.TypeObject,
				Description: "Search basic Wikipedia.",
				Properties: map[string]*genai.Schema{
					"usetool": {Type: genai.TypeBoolean, Description: "True if Wikipedia search is needed."},
					"query":   {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "Independent search queries."},
				},
			},
			"deepwikipediatool": {
				Type:        genai.TypeObject,
				Description: "Search deep/related Wikipedia pages.",
				Properties: map[string]*genai.Schema{
					"usetool": {Type: genai.TypeBoolean, Description: "True if deep Wikipedia search is needed."},
					"query":   {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "Independent search queries."},
				},
			},
			"news_searchtool": {
				Type:        genai.TypeObject,
				Description: "Search for recent news articles.",
				Properties: map[string]*genai.Schema{
					"usetool": {Type: genai.TypeBoolean, Description: "True if news search is needed."},
					"query":   {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "Independent search queries."},
				},
			},
			"github_searchtool": {
				Type:        genai.TypeObject,
				Description: "Search GitHub repositories or code.",
				Properties: map[string]*genai.Schema{
					"usetool": {Type: genai.TypeBoolean, Description: "True if GitHub search is needed."},
					"query":   {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "Independent search queries."},
				},
			},
			"stackoverflow_searchtool": {
				Type:        genai.TypeObject,
				Description: "Search programming questions on Stack Overflow.",
				Properties: map[string]*genai.Schema{
					"usetool": {Type: genai.TypeBoolean, Description: "True if Stack Overflow search is needed."},
					"query":   {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "Independent search queries."},
				},
			},
			"google_search_aitool": {
				Type:        genai.TypeObject,
				Description: "AI-enhanced Google search for better page summaries.",
				Properties: map[string]*genai.Schema{
					"usetool": {Type: genai.TypeBoolean, Description: "True if AI-powered Google search is needed."},
					"query":   {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "Independent search queries."},
				},
			},
			"google_searchtool": {
				Type:        genai.TypeObject,
				Description: "Standard Google search.",
				Properties: map[string]*genai.Schema{
					"usetool": {Type: genai.TypeBoolean, Description: "True if Google search is needed."},
					"query":   {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "Independent search queries."},
				},
			},
			"youtube_playlisttool": {
				Type:        genai.TypeObject,
				Description: "Search YouTube playlists.",
				Properties: map[string]*genai.Schema{
					"usetool": {Type: genai.TypeBoolean, Description: "True if YouTube playlist search is needed."},
					"query":   {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "Independent search queries."},
				},
			},
			"youtube_videotool": {
				Type:        genai.TypeObject,
				Description: "Search YouTube videos.",
				Properties: map[string]*genai.Schema{
					"usetool": {Type: genai.TypeBoolean, Description: "True if YouTube video search is needed."},
					"query":   {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "Independent search queries."},
				},
			},
			"weathertool": {
				Type:        genai.TypeObject,
				Description: "Get current weather for places.",
				Properties: map[string]*genai.Schema{
					"usetool": {Type: genai.TypeBoolean, Description: "True if weather data is needed."},
					"query":   {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "Array of place names only."},
				},
			},
			"reddit_searchtool": {
				Type:        genai.TypeObject,
				Description: "Search Reddit for discussions.",
				Properties: map[string]*genai.Schema{
					"usetool": {Type: genai.TypeBoolean, Description: "True if Reddit search is needed."},
					"query":   {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "Independent search queries."},
				},
			},
			"usesitecraft": {
				Type:        genai.TypeObject,
				Description: "it makes website/web app using. only use this tool when user ask for website",
				Properties: map[string]*genai.Schema{
					"usetool": {Type: genai.TypeBoolean, Description: "True if want to make website."},
					"query":   {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "queries for sitecraft to make website"},
				},
			},
		},
		Required: []string{
			"wikipediatool", "deepwikipediatool", "news_searchtool",
			"github_searchtool", "stackoverflow_searchtool",
			"google_search_aitool", "google_searchtool",
			"youtube_playlisttool", "youtube_videotool",
			"weathertool", "reddit_searchtool",
		},
	},
}

var GenerateCode = genai.FunctionDeclaration{
	Name:        "GenerateCode",
	Description: "Generates and returns valid and working code for particular file for nextjs project. The code must be self-contained, properly typed, and not include explanations or extra text.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"code": {
				Type:        genai.TypeString,
				Description: "A complete, valid code for the specified file, including all necessary imports and exports.",
			},
		},
		Required: []string{"code"},
	},
}

var FilenameTool = genai.FunctionDeclaration{
	Name:        "file",
	Description: "Generates and returns valid Next.js project files as an array, including frontend and backend. Use pages/{page}.js for frontend pages and pages/api/{endpoint}.js for backend endpoints. Only return file names with extensions; do not include explanations.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"frontendFiles": {
				Type:        genai.TypeArray,
				Description: "Frontend page files with extensions (.js or .jsx), e.g., pages/index.js, pages/about.jsx",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
			"backendFiles": {
				Type:        genai.TypeArray,
				Description: "Backend API files (.js) inside pages/api/, e.g., pages/api/hello.js, pages/api/data.js",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"frontendFiles", "backendFiles"},
	},
}

var PostCode = genai.FunctionDeclaration{
	Name:        "files",
	Description: "",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"frontendCode": {
				Type:        genai.TypeObject,
				Description: "A map where keys are NextJS frontent component filenames (App.js or index.js or etc) dont include path and values are their fixed, valid code.",
			},
			"backendCode": {
				Type:        genai.TypeObject,
				Description: "A map where keys are NextJS Backend component filenames (App.js or index.js or etc) dont include path it is assumed it inside api and values are their fixed, valid code.",
			},
		},
		Required: []string{"frontendCode", "backendCode"},
	},
}

var RAG = genai.FunctionDeclaration{
	Name:        "rag",
	Description: "Retrieves reusable Tailwind components for building websites",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"rag": {
				Type:        genai.TypeArray,
				Description: "An array of desgin styles to search on vector db",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
			"img": {
				Type:        genai.TypeArray,
				Description: "An array of image search queries that requried for website to make it cool max 1 image query",
				Items:       &genai.Schema{Type: genai.TypeString},
			},
		},
		Required: []string{"rag", "img"},
	},
}
