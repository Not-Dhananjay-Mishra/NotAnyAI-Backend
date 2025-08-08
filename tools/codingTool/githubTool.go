package codingtool

import (
	"context"
	"fmt"
	"log"
	"server/utils"

	"github.com/google/go-github/v74/github"
)

func GetReadMeGithub(owner string, repo string) string {
	client := github.NewClient(nil).WithAuthToken(utils.GITHUB_API)

	ctx := context.Background()
	readme, _, err := client.Repositories.GetReadme(ctx, owner, repo, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	content, err := readme.GetContent()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return content
}

func GetRepoGithub(query string) []utils.GithubRepo {
	var Repositories []utils.GithubRepo
	client := github.NewClient(nil).WithAuthToken(utils.GITHUB_API)

	opts := &github.SearchOptions{
		Sort:  "stars",
		Order: "desc",
	}

	result, _, err := client.Search.Repositories(context.Background(), query, opts)
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}
	for _, repo := range result.Repositories {
		Repositories = append(Repositories, utils.GithubRepo{FullName: *repo.FullName, RepoName: *repo.Name, OwnerName: *repo.Owner.Login, URL: *repo.URL})
	}
	return Repositories
}

func SearchInRepo(owner string, repo string, query string) {
	client := github.NewClient(nil).WithAuthToken(utils.GITHUB_API)
	ctx := context.Background()

	searchQuery := fmt.Sprintf("%s repo:%s/%s", query, owner, repo)
	opts := &github.SearchOptions{Sort: "indexed", Order: "desc"}

	result, _, err := client.Search.Code(ctx, searchQuery, opts)
	if err != nil {
		log.Fatalf("Search error: %v", err)
	}

	fmt.Println(*result.CodeResults[0].Path)
}
