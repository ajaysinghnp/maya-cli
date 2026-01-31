package resolver

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
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
	log.Success("MoviesBazar webpage requested successfully!")
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error(fmt.Sprintf("Unexpected HTTP status: %d", resp.StatusCode))
		return nil, errors.New("failed to fetch page")
	}
	log.Success("MoviesBazar webpage fetched successfully!")

	// 2️⃣ Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Error("Failed to parse HTML: " + err.Error())
		return nil, err
	}
	log.Success("MoviesBazar webpage parsed successfully!")

	// 3️⃣ Extract metadata
	meta := &metadata.Metadata{
		Type: metadata.Movie, // default to movie, could detect later
	}

	// 3️⃣ Extract all fields using helpers
	meta.Title = extractTitle(doc, log)
	meta.Year = extractYear(doc, log)
	meta.ReleaseDate = extractReleaseDate(doc, log)
	meta.Language = extractLanguage(doc, log)
	meta.Cast = extractCast(doc, log)
	meta.Genres = extractGenre(doc, log)
	meta.IDs = extractIDs(doc, log)
	meta.M3U8URLs = extractM3U8URLs(doc, log)

	log.Success(fmt.Sprintf("Fetched metadata for '%s' (%s)", meta.Title, meta.Year))
	return meta, nil
}

///// Helper functions /////

// Generic prompt for user input
func promptInput(prompt string, log iface.Logger) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Warn("Failed to read input, leaving empty")
		return ""
	}
	return strings.TrimSpace(input)
}

// Extract movie title
func extractTitle(doc *goquery.Document, log iface.Logger) string {
	title := strings.TrimSpace(doc.Find(".mobile\\:px-2\\.5 > div:nth-child(1) strong:contains('Title:') + h1").Text())
	if title == "" {
		log.Warn("Title not found, prompting user")
		title = promptInput("Enter movie title: ", log)
	}
	log.Success("Title: " + title)
	return title
}

// Extract release year
func extractYear(doc *goquery.Document, log iface.Logger) string {
	year := strings.TrimSpace(doc.Find(".mobile\\:px-2\\.5 > div:nth-child(1) strong:contains('Year:') + div").Text())
	if year == "" {
		log.Warn("Year not found, prompting user")
		year = promptInput("Enter year: ", log)
	}
	log.Success("Year: " + year)
	return year
}

// Extract release date
func extractReleaseDate(doc *goquery.Document, log iface.Logger) string {
	date := strings.TrimSpace(doc.Find(".mobile\\:px-2\\.5 > div:nth-child(1) strong:contains('Released:') + div").Text())
	if date == "" {
		log.Warn("Release date not found, optional")
		date = promptInput("Enter release date (optional): ", log)
	}
	log.Success("Release Date: " + date)
	return date
}

// Extract language
func extractLanguage(doc *goquery.Document, log iface.Logger) string {
	lang := strings.TrimSpace(doc.Find(".mobile\\:px-2\\.5 > div:nth-child(1) strong:contains('Language:') + a").Text())
	log.Success("Language: " + lang)
	return lang
}

// Extract cast as []string
func extractCast(doc *goquery.Document, log iface.Logger) []string {
	castList := []string{}
	doc.Find(".mobile\\:px-2\\.5 > div:nth-child(1) strong:contains('Star cast:') ~ div").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		text = strings.TrimSuffix(text, ",") // remove trailing comma
		if text != "" {
			castList = append(castList, text)
		}
	})
	if len(castList) > 0 {
		log.Success("Cast: " + strings.Join(castList, ", "))
	} else {
		log.Warn("No cast found")
	}
	return castList
}

// Extract genre as []string
func extractGenre(doc *goquery.Document, log iface.Logger) []string {
	genres := []string{}
	doc.Find(".mobile\\:px-2\\.5 > div:nth-child(1) strong:contains('Genre:') + a").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			genres = append(genres, text)
		}
	})
	if len(genres) > 0 {
		log.Success("Genres: " + strings.Join(genres, " • "))
	} else {
		log.Warn("No genres found")
	}
	return genres
}

// extractIDs tries to extract TMDB or IMDB IDs from the movie poster img src
func extractIDs(doc *goquery.Document, log iface.Logger) metadata.IDs {
	ids := metadata.IDs{}

	// Use the container that has the img poster
	img := doc.Find("div.w-full.max-w-[300px].max-h-[400px] img")
	if img.Length() == 0 {
		log.Warn("Poster image not found, cannot infer TMDB/IMDB IDs")
		return ids
	}

	src, exists := img.Attr("src")
	if !exists || src == "" {
		log.Warn("Poster image src not found, cannot infer TMDB/IMDB IDs")
		return ids
	}

	// Example src formats:
	// https://image.tmdb.org/t/p/w500/abc123.jpg
	// https://www.imdb.com/title/tt1234567/mediaviewer/rm123456789

	if strings.Contains(src, "tmdb.org") {
		parts := strings.Split(src, "/")
		for i, part := range parts {
			if part == "movie" && i+1 < len(parts) {
				ids.TMDB = parts[i+1]
				log.Success("Detected TMDB ID: " + ids.TMDB)
				break
			}
		}
	} else if strings.Contains(src, "imdb.com") {
		parts := strings.Split(src, "/")
		for _, part := range parts {
			if strings.HasPrefix(part, "tt") { // imdb IDs start with tt
				ids.IMDB = part
				log.Success("Detected IMDB ID: " + ids.IMDB)
				break
			}
		}
	} else {
		log.Warn("Could not detect TMDB or IMDB ID from poster src")
	}

	return ids
}

// Extract M3U8 URLs from script/src/href
func extractM3U8URLs(doc *goquery.Document, log iface.Logger) []string {
	urls := []string{}
	doc.Find("script, a, source").Each(func(i int, s *goquery.Selection) {
		for _, attr := range []string{"src", "href"} {
			if val, exists := s.Attr(attr); exists && strings.HasSuffix(val, ".m3u8") {
				urls = append(urls, val)
			}
		}
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
	if len(urls) > 0 {
		log.Success(fmt.Sprintf("Found %d M3U8 URL(s)", len(urls)))
	} else {
		log.Warn("No M3U8 URLs found")
	}
	return urls
}
