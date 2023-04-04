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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/jaytaylor/html2text"
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

func isAnyLinkMatch(absoluteLink string) bool {
	for _, pattern := range LinkPatterns {
		if pattern.MatchString(absoluteLink) {
			return true
		}
	}
	return false
}

func shouldContinueCrawl(text, absoluteLink string) bool {
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
	return isAnyLinkMatch(absoluteLink)

}

func shouldFetchJobListing(text, absoluteLink string) bool {
	return isAnyLinkMatch(absoluteLink)
}

type JobPosting struct {
	AbsoluteURL    string
	ListingContent string
	CompanyName    string
	JobBoard       string
	Md5Hash        string
	FetchedAt      time.Time
}

func processHtml(e *colly.HTMLElement) (JobPosting, error) {
	hostname := e.Request.URL.Host
	absoluteUrl := e.Request.URL.String()
	listingContent, err := html2text.FromString(e.Text,
		html2text.Options{
			PrettyTables: false,
			OmitLinks:    true,
			TextOnly:     true,
		})
	if err != nil {
		log.Printf("error converting html to text: %v", err)
		return JobPosting{}, err
	}
	hash := md5.Sum([]byte(listingContent))
	// Convert the hash to a hex string
	md5Hash := hex.EncodeToString(hash[:])
	log.Printf("Hashed body at link: %q \n", e.Request.URL.Host)

	var jobBoard string
	var companyName string

	if strings.Contains(e.Request.URL.Host, "jobs.lever") {
		jobBoard = "lever"
		companyName = strings.Split(absoluteUrl, "/")[3]

	} else if strings.Contains(e.Request.URL.Host, "ycombinator") {
		jobBoard = "ycombinator"
		companyName = strings.Split(absoluteUrl, "/")[4]

	} else if strings.Contains(e.Request.URL.Host, "greenhouse") {
		jobBoard = "greenhouse"
		companyName = strings.Split(absoluteUrl, "/")[3]
	} else {
		return JobPosting{}, fmt.Errorf("unhandled hostname: %s", hostname)
	}

	return JobPosting{
		absoluteUrl,
		listingContent,
		companyName,
		jobBoard,
		md5Hash,
		time.Now(),
	}, nil
}

type JobListingFetcher struct {
	collector    *colly.Collector
	htmlSelector string
}

var Fetchers = map[string]JobListingFetcher{
	"ycombinator": {
		colly.NewCollector(colly.DisallowedDomains("account.ycombinator.com")),
		"div.mx-auto.max-w-ycdc-page > section > div > div.flex-grow.space-y-5"},
	"greenhouse": {colly.NewCollector(colly.Async(true)), "#content"},
	"jobs.lever": {colly.NewCollector(colly.Async(true)), "body > div.content-wrapper.posting-page > div"},
}

func initCrawlers(messages chan<- JobPosting) []*colly.Collector {
	// Instantiate default collector
	ycCrawler := colly.NewCollector(
		colly.MaxDepth(5),
		colly.Async(true),
	)
	// Set max Parallelism and introduce a Random Delay
	ycCrawler.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	// On every a element which has href attribute call callback
	ycCrawler.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		absoluteLink := e.Request.AbsoluteURL(link)

		if shouldContinueCrawl(e.Text, absoluteLink) {
			ycCrawler.Visit(absoluteLink)
		}
		if shouldFetchJobListing(e.Text, absoluteLink) {
			for host, fetcher := range Fetchers {
				if strings.Contains(e.Request.URL.Host, host) {
					fetcher.collector.Visit(absoluteLink)
					break
				}
			}
		}

	})
	ycCrawler.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	for host, fetcher := range Fetchers {
		fetcher.collector.OnHTML(fetcher.htmlSelector, func(e *colly.HTMLElement) {
			posting, err := processHtml(e)
			if err != nil {
				log.Printf("failed to parse %s page %e", host, err)
			}
			messages <- posting
		})
	}

	// Starting point
	ycCrawler.Visit("https://news.ycombinator.com/jobs")

	// Aggregate colly collectors so we can .Wait on all of them
	collectors := []*colly.Collector{}
	collectors = append(collectors, ycCrawler)
	for _, fetcher := range Fetchers {
		collectors = append(collectors, fetcher.collector)
	}
	return collectors
}

func Run() {
	messages := make(chan JobPosting)
	crawler := initCrawlers(messages)
	var count int

	go func(messages <-chan JobPosting) {
		for msg := range messages {
			log.Printf("Received Job Posting %s", msg)
			count++
		}
	}(messages)

	for _, c := range crawler {
		c.Wait()
	}
	close(messages)
	log.Printf("Found %d jobs", count)
}
