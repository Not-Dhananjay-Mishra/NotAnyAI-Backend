package individualtool

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"server/utils"
)

type YouTubeResponse struct {
	Items []Item `json:"items"`
}

type Item struct {
	ID VideoID `json:"id"`
}

type VideoID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

type PYouTubeResponse struct {
	Items []PItem `json:"items"`
}

type PItem struct {
	ID PlaylistID `json:"id"`
}

type PlaylistID struct {
	Kind       string `json:"kind"`
	PlaylistID string `json:"playlistId"`
}
type PlaylistResponse struct {
	Items []PlaylistItem `json:"items"`
}

type PlaylistItem struct {
	Snippet        PlaylistSnippet        `json:"snippet"`
	ContentDetails PlaylistContentDetails `json:"contentDetails"`
}

type PlaylistSnippet struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PlaylistContentDetails struct {
	ItemCount int `json:"itemCount"`
}
type VideoResponse struct {
	Items []VideoItem `json:"items"`
}

type VideoItem struct {
	Snippet        VideoSnippet        `json:"snippet"`
	ContentDetails VideoContentDetails `json:"contentDetails"`
	Statistics     VideoStatistics     `json:"statistics"`
}

type VideoSnippet struct {
	PublishedAt  string `json:"publishedAt"`
	ChannelID    string `json:"channelId"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	ChannelTitle string `json:"channelTitle"`
}

type VideoContentDetails struct {
	Duration string `json:"duration"`
}

type VideoStatistics struct {
	ViewCount    string `json:"viewCount"`
	LikeCount    string `json:"likeCount"`
	CommentCount string `json:"commentCount"`
}

func YoutubeToolVideo(query string) []utils.VideoDetails {
	var Videos []utils.VideoDetails
	baseURL := "https://www.googleapis.com/youtube/v3/search"
	params := url.Values{}
	params.Add("key", utils.YOUTUBE_API)
	params.Add("maxResults", "7")
	params.Add("type", "video")
	params.Add("q", query)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	var youtubeResponse YouTubeResponse
	if err := json.Unmarshal(bodyBytes, &youtubeResponse); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return nil
	}
	for _, item := range youtubeResponse.Items {
		if item.ID.Kind == "youtube#video" {
			title, channel, view, like := GetDetailsVideo(item.ID.VideoID)
			Videos = append(Videos, utils.VideoDetails{Link: fmt.Sprintf("https://www.youtube.com/watch?v=%s", item.ID.VideoID), Title: title, ChannelName: channel, View: view, Like: like})
		}
	}
	return Videos
}
func YoutubeToolPlaylist(query string) ([]string, []string, []string) {
	var Links []string
	var Titles []string
	var Descriptions []string
	baseURL := "https://www.googleapis.com/youtube/v3/search"
	params := url.Values{}
	params.Add("key", utils.YOUTUBE_API)
	params.Add("maxResults", "5")
	params.Add("type", "playlist")
	params.Add("q", query)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return nil, nil, nil
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, nil, nil
	}

	var youtubeResponse PYouTubeResponse
	if err := json.Unmarshal(bodyBytes, &youtubeResponse); err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return nil, nil, nil
	}

	for _, item := range youtubeResponse.Items {
		if item.ID.Kind == "youtube#playlist" {
			playlistURL := fmt.Sprintf("https://www.youtube.com/playlist?list=%s", item.ID.PlaylistID)
			Links = append(Links, playlistURL)
			title, desc, _ := GetDetailsPlaylist(item.ID.PlaylistID)
			Titles = append(Titles, title)
			Descriptions = append(Descriptions, desc)
		}
	}
	return Links, Titles, Descriptions

}

func GetDetailsPlaylist(id string) (string, string, int) {
	baseURL := "https://www.googleapis.com/youtube/v3/playlists"
	params := url.Values{}
	params.Add("key", utils.YOUTUBE_API)
	params.Add("part", "snippet,contentDetails")
	params.Add("id", id)
	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return "", "", 0
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", "", 0
	}

	var playlistResponse PlaylistResponse
	if err := json.Unmarshal(bodyBytes, &playlistResponse); err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return "", "", 0
	}

	if len(playlistResponse.Items) > 0 {
		playlist := playlistResponse.Items[0]
		return playlist.Snippet.Title, playlist.Snippet.Description, playlist.ContentDetails.ItemCount
	} else {
		return "", "", 0
	}
}
func GetDetailsVideo(id string) (string, string, string, string) {
	baseURL := "https://www.googleapis.com/youtube/v3/videos"
	params := url.Values{}
	params.Add("key", utils.YOUTUBE_API)
	params.Add("part", "snippet,contentDetails,statistics")
	params.Add("id", id)

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		log.Printf("Error making request: %v", err)
		return "", "", "", ""
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return "", "", "", ""
	}

	var videoResponse VideoResponse
	if err := json.Unmarshal(bodyBytes, &videoResponse); err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return "", "", "", ""
	}

	if len(videoResponse.Items) > 0 {
		video := videoResponse.Items[0]
		return video.Snippet.Title, video.Snippet.ChannelTitle, video.Statistics.ViewCount, video.Statistics.LikeCount
	} else {
		return "", "", "", ""
	}
}
