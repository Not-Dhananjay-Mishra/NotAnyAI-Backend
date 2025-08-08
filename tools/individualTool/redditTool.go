package individualtool

import (
	"context"
	"encoding/json"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func RedditTool(subreddit string, query string) (string, []string) {
	var title string
	var comment []string

	client, _ := reddit.NewReadonlyClient()
	opts := &reddit.ListSubredditOptions{
		ListOptions: reddit.ListOptions{Limit: 5},
		Sort:        "year",
	}
	popts := &reddit.ListPostSearchOptions{ListPostOptions: reddit.ListPostOptions{ListOptions: reddit.ListOptions{Limit: 1}}}
	sub, _, err := client.Subreddit.Search(context.Background(), subreddit, opts)
	if err != nil {
		return "nil", nil
	}
	//fmt.Printf(sub[0].Name)
	posts, _, err := client.Subreddit.SearchPosts(context.Background(), query, sub[0].Name, popts)
	if err != nil {
		return "nil", nil
	}
	title = posts[0].Title
	//fmt.Println(posts[0].Body)
	comments, _, err := client.Post.Get(context.Background(), posts[0].ID)
	if err != nil {
		return "nil", nil
	}
	for i := range comments.Comments {
		res, _ := json.Marshal(comments.Comments[i].Body)
		comment = append(comment, string(res))
	}
	return title, comment
}
