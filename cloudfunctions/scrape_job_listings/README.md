# Scrape Job Listings

Crawlers here are responsible for fetching job listings from known sources. Currently we are handling
- lever
- ycombinator
- greenhouse

Crawlers use heuristics to match link urls and product messages for further classification and information extraction.

For information extraction we use document selection rules. These are brittle but simple.

