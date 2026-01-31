package resolver

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery" // for HTML parsing
	"github.com/ajaysinghnp/maya-cli/internal/logger/iface"
	"github.com/ajaysinghnp/maya-cli/internal/metadata"
)

// fetchMoviesBazarMetadata fetches metadata and m3u8 URLs from a MoviesBazar page
func fetchMoviesBazarMetadata(url string, log iface.Logger) (*metadata.Metadata, error) {
	log.Info("Fetching metadata from MoviesBazar for URL: " + url)

	// 1️⃣ Download the page
	resp, err := http.Get(url)
	if err != nil {
		log.Error("Failed to download page: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error(fmt.Sprintf("Unexpected HTTP status: %d", resp.StatusCode))
		return nil, errors.New("failed to fetch page")
	}

	// 2️⃣ Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error("Failed to parse HTML: " + err.Error())
		return nil, err
	}

	// 3️⃣ Extract metadata
	meta := &metadata.Metadata{
		Type: metadata.Movie, // default to movie, could detect later
	}

	// Example: extract title
	meta.Title = strings.TrimSpace(doc.Find("h1.movie-title").Text())
	if meta.Title == "" {
		log.Warn("Could not detect movie title, prompting user")
		return metadata.PromptUser(log)
	}

	// Example: extract year
	yearText := strings.TrimSpace(doc.Find("span.year").Text())
	meta.Year = yearText

	// Example: extract TMDB/IMDB IDs if present in the page
	meta.IDs = metadata.IDs{
		IMDB: extractIMDB(doc),
		TMDB: extractTMDB(doc),
	}

	// 4️⃣ Extract M3U8 URLs from page
	meta.M3U8URLs = extractM3U8URLs(doc, log)
	if len(meta.M3U8URLs) == 0 {
		log.Warn("No M3U8 URLs found, download may fail")
	}

	log.Success(fmt.Sprintf("Fetched metadata for '%s' (%s)", meta.Title, meta.Year))
	return meta, nil
}

// extractIMDB and extractTMDB are helpers to find IDs in page
func extractIMDB(doc *goquery.Document) string {
	// Example: find link with imdb.com/title/ttXXXXX
	link := ""
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(href, "imdb.com/title/") {
			link = strings.TrimSpace(strings.Split(href, "?")[0])
		}
	})
	return link
}

func extractTMDB(doc *goquery.Document) string {
	// Example: find link with themoviedb.org/movie/XXXXX
	link := ""
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(href, "themoviedb.org/movie/") {
			link = strings.TrimSpace(strings.Split(href, "?")[0])
		}
	})
	return link
}

// extractM3U8URLs searches for .m3u8 URLs in the page's HTML
func extractM3U8URLs(doc *goquery.Document, log iface.Logger) []string {
	urls := []string{}
	doc.Find("script, a, source").Each(func(i int, s *goquery.Selection) {
		// check src/href/content for .m3u8
		for _, attr := range []string{"src", "href"} {
			if val, exists := s.Attr(attr); exists && strings.HasSuffix(val, ".m3u8") {
				urls = append(urls, val)
			}
		}

		// Also check inline script contents
		content := s.Text()
		if strings.Contains(content, ".m3u8") {
			parts := strings.Split(content, "\"")
			for _, p := range parts {
				if strings.HasSuffix(p, ".m3u8") {
					urls = append(urls, p)
				}
			}
		}
	})
	return urls
}
