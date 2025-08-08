package individualtool

import (
	"fmt"

	gowiki "github.com/trietmn/go-wiki"
)

func WikiSummarySearch(query string) string {
	res, err := gowiki.Summary(query, 10, -1, false, true)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func WikiDeepSearch(query string) string {
	// Get the page
	page, err := gowiki.GetPage(query, -1, false, true)
	if err != nil {
		fmt.Println(err)
	}

	// Get the content of the page
	content, err := page.GetContent()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("This is the page content: %v\n", content)
	return content
}
