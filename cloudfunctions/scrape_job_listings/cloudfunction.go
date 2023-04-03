/*
Perform relatively shallow crawl, using heuristics to classify JobListings and traversal continuations.

There are two colly collectors
1. Traversing links starting from the hn job board
2. Collecting HTML bodies and their hashes for further processing

For the most part, the heuristics look a links
- Yc Jobs
- Lever
- Greenhouse
*/
package cloudfunction

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

var YcPattern = regexp.MustCompile(`.*(ycombinator\.com)\/companies.*\/jobs\/`)

// handling https://jobs.lever.co/memfault/730541eb-637f-4d9d-9526-8949432f9a34
// to check for UUIDs
var LeverPattern = regexp.MustCompile(`.*jobs\.lever\.co\/[a-zA-Z0-9._%+-]+\/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// handling https://boards.greenhouse.io/supabase/jobs/4777008004
var GreenhousePattern = regexp.MustCompile(`.*(boards\.greenhouse\.io)\/.*\/jobs\/`)

var LinkPatterns = []*regexp.Regexp{
	YcPattern,
	LeverPattern,
	GreenhousePattern,
}

func CheckLinkPatterns(absoluteLink string) bool {
	for _, pattern := range LinkPatterns {
		if pattern.MatchString(absoluteLink) {
			return true
		}
	}
	return false
}

func ShouldContinueCrawl(text, absoluteLink string) bool {
	searchStr := strings.ToLower(text)
	// ends with /careers
	// ends with /jobs
	if strings.HasSuffix(searchStr, "jobs") || strings.HasSuffix(searchStr, "careers") {
		return true
	}
	// looks like a hn posting
	if strings.Contains(searchStr, "hiring") {
		return true
	}
	return CheckLinkPatterns(absoluteLink)

}

func ShouldFetchJobListing(text, absoluteLink string) bool {
	return CheckLinkPatterns(absoluteLink)
}

func InitCrawler() *colly.Collector {
	// Instantiate default collector
	ycCrawler := colly.NewCollector(
		colly.MaxDepth(3),
		colly.Async(true),
	)

	jobPostingFetch := colly.NewCollector(
		colly.Async(false),
		colly.Debugger(&debug.LogDebugger{}),
		colly.DisallowedDomains("account.ycombinator.com"),
	)

	jobPostingFetch.OnResponse(func(r *colly.Response) {
		// Hash the body of the response using MD5
		// hash := md5.Sum(r.Body)
		// text := string(r.Body)

		// Convert the hash to a hex string
		// hashString := hex.EncodeToString(hash[:])
		log.Printf("Hashed body at link: %q \n", r.Request.URL)

	})
	// Set max Parallelism and introduce a Random Delay
	ycCrawler.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	// On every a element which has href attribute call callback
	ycCrawler.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteLink := e.Request.AbsoluteURL(link)
		// log.Printf("Link found: %q -> %s\n", e.Text, absoluteLink)

		if ShouldContinueCrawl(e.Text, absoluteLink) {
			ycCrawler.Visit(absoluteLink)
		}
		if ShouldFetchJobListing(e.Text, absoluteLink) {
			jobPostingFetch.Visit(absoluteLink)
		}

	})
	ycCrawler.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	// Starting point
	ycCrawler.Visit("https://news.ycombinator.com/jobs")
	return ycCrawler
}

func Run() {
	crawler := InitCrawler()
	crawler.Wait()
}
