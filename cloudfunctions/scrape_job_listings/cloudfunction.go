/*
Perform relatively shallow crawl, using heuristics to classify JobListings and traversal continuations.

There are two colly collectors
1. Traversing links starting from the hn job board
2. Collecting HTML bodies and their hashes for further processing
*/
package cloudfunction

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/getsentry/sentry-go"

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
	// looks like a hn listing
	if strings.Contains(searchStr, "hiring") {
		return true
	}
	return isAnyLinkMatch(absoluteLink)

}

func shouldFetchJobListing(text, absoluteLink string) bool {
	return isAnyLinkMatch(absoluteLink)
}

type JobListing struct {
	AbsoluteURL    string
	ListingContent string
	CompanySlug    string
	JobBoard       string
	Md5Hash        string
	FetchedAt      time.Time
}

func processHtml(e *colly.HTMLElement) (JobListing, error) {
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
		return JobListing{}, err
	}
	hash := md5.Sum([]byte(listingContent))
	// Convert the hash to a hex string
	md5Hash := hex.EncodeToString(hash[:])
	log.Printf("Hashed body at link: %q \n", e.Request.URL.Host)

	var jobBoard string
	var companySlug string

	if strings.Contains(e.Request.URL.Host, "jobs.lever") {
		jobBoard = "lever"
		companySlug = strings.Split(absoluteUrl, "/")[3]

	} else if strings.Contains(e.Request.URL.Host, "ycombinator") {
		jobBoard = "ycombinator"
		companySlug = strings.Split(absoluteUrl, "/")[4]

	} else if strings.Contains(e.Request.URL.Host, "greenhouse") {
		jobBoard = "greenhouse"
		companySlug = strings.Split(absoluteUrl, "/")[3]
	} else {
		return JobListing{}, fmt.Errorf("unhandled hostname: %s", hostname)
	}

	return JobListing{
		absoluteUrl,
		listingContent,
		companySlug,
		jobBoard,
		md5Hash,
		time.Now(),
	}, nil
}

type JobListingFetcher struct {
	collector    *colly.Collector
	htmlSelector string
}

var fetchers = map[string]JobListingFetcher{
	"ycombinator": {
		colly.NewCollector(colly.DisallowedDomains("account.ycombinator.com")),
		"div.mx-auto.max-w-ycdc-page > section > div > div.flex-grow.space-y-5"},
	"greenhouse": {colly.NewCollector(colly.Async(false)), "#content"},
	"jobs.lever": {colly.NewCollector(colly.Async(false)), "body > div.content-wrapper.posting-page > div"},
}

func initCrawlers(messages chan<- JobListing) *colly.Collector {
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
			for host, fetcher := range fetchers {
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

	for host, fetcher := range fetchers {
		fetcher.collector.OnHTML(fetcher.htmlSelector, func(e *colly.HTMLElement) {
			listing, err := processHtml(e)
			if err != nil {
				log.Printf("failed to parse %s page %e", host, err)
			}
			messages <- listing
		})
	}

	// Starting point
	ycCrawler.Visit("https://news.ycombinator.com/jobs")

	// Aggregate colly collectors so we can .Wait on all of them
	return ycCrawler
}

func Run() {
	messages := make(chan JobListing)
	crawler := initCrawlers(messages)
	var count int

	go func(messages <-chan JobListing) {
		for msg := range messages {
			log.Println("Received Job Listing. Not sending message.", msg)
			count++
		}
	}(messages)

	crawler.Wait()
	close(messages)
	log.Printf("Completed. Found %d jobs", count)
}

func init() {
	functions.HTTP("Handler", handler)
}

func handleError(w http.ResponseWriter, msg string, err error) {
	err = fmt.Errorf("%s: %w", msg, err)
	log.Print(err)
	sentry.CaptureException(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)

}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		ServerName:       os.Getenv("FUNCTION_NAME"),
	})
	if err != nil {
		handleError(w, "Error intializing sentry", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.RecoverWithContext(ctx)
	Run()
}
