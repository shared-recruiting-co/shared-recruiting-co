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

type JobPosting struct {
	AbsoluteURL    string
	ListingContent string
	CompanyName    string
	JobBoard       string
	Md5Hash        string
	FetchedAt      time.Time
}

func ProcessHtml(e *colly.HTMLElement) (JobPosting, error) {
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
	log.Printf("Hashed body at link: %q \n", e.Request.URL.Host)
	hash := md5.Sum([]byte(listingContent))
	// Convert the hash to a hex string
	md5Hash := hex.EncodeToString(hash[:])

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

func InitCrawler() *colly.Collector {
	// Instantiate default collector
	ycCrawler := colly.NewCollector(
		colly.MaxDepth(3),
		colly.Async(true),
	)
	// Set max Parallelism and introduce a Random Delay
	ycCrawler.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	ycFetcher := colly.NewCollector(
		colly.DisallowedDomains("account.ycombinator.com"),
	)
	greenhouseFetcher := colly.NewCollector()
	leverFetcher := colly.NewCollector()

	ycFetcher.OnHTML("div.mx-auto.max-w-ycdc-page > section > div > div.flex-grow.space-y-5", func(e *colly.HTMLElement) {
		log.Println(ProcessHtml(e))
	})

	greenhouseFetcher.OnHTML("#content", func(e *colly.HTMLElement) {
		log.Println(ProcessHtml(e))
	})

	leverFetcher.OnHTML("body > div.content-wrapper.posting-page > div", func(e *colly.HTMLElement) {
		log.Println(ProcessHtml(e))

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
			if strings.Contains(e.Request.URL.Host, "jobs.lever") {
				leverFetcher.Visit(absoluteLink)

			} else if strings.Contains(e.Request.URL.Host, "ycombinator") {
				ycFetcher.Visit(absoluteLink)

			} else if strings.Contains(e.Request.URL.Host, "greenhouse") {
				greenhouseFetcher.Visit(absoluteLink)
			}
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
