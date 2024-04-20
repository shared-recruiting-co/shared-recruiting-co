package main

import (
	scrape "github.com/shared-recruiting-co/shared-recruiting-co/cloudfunctions/scrape_job_listings"
)

func main() {
	scrape.Run()
}
