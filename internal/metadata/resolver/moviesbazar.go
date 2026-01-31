package resolver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/PuerkitoBio/goquery"
	"github.com/ajaysinghnp/maya-cli/internal/logger"
	"github.com/ajaysinghnp/maya-cli/internal/logger/iface"
	"github.com/ajaysinghnp/maya-cli/internal/metadata"
	"github.com/ajaysinghnp/maya-cli/utils"
)

// fetchMoviesBazarMetadata fetches metadata and m3u8 URLs from a MoviesBazar page
func fetchMoviesBazarMetadata(url string, log iface.Logger) (*metadata.Metadata, error) {
	log.Info("Fetching metadata from MoviesBazar for URL: " + url)

	resp, err := http.Get(url)
	if err != nil {
		log.Error("Failed to download page: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	log.Success("Movie info requested successfully!")

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	log.Success("Movie page fetched successfully!")

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Success("Movie page parsed successfully!")

	scriptData, err := extractMovieDetailsFromScript(doc, log)
	if err != nil {
		return nil, err
	}

	meta := &metadata.Metadata{
		Title:           scriptData.Title,
		Year:            utils.NormalizeYear(scriptData.Year),
		Category:        scriptData.Category,
		Type:            metadata.Movie,
		ReleaseDate:     utils.NormalizeDate(scriptData.ReleaseDate),
		Language:        utils.NormalizeLanguage(scriptData.Language),
		Genres:          utils.NormalizeSlice(scriptData.Genres),
		Cast:            utils.NormalizeSlice(scriptData.Cast),
		Thumbnail:       scriptData.Thumbnail,
		HlsSourceDomain: scriptData.HlsSourceDomain,
		Sources:         scriptData.Sources,
		IDs: metadata.IDs{
			IMDB: scriptData.IDs.IMDB,
			TMDB: scriptData.IDs.TMDB,
		},
	}

	log.Success(fmt.Sprintf(
		"Fetched metadata for '%s' (%d)",
		meta.Title,
		meta.Year,
	))

	return meta, nil
}

///// Script Parsing /////

// extractMovieDetailsFromScript parses self.__next_f.push([...]) blocks
func extractMovieDetailsFromScript(doc *goquery.Document, log iface.Logger) (*metadata.Metadata, error) {
	log.Info("Initiating the extraction of movie info...")

	re := regexp.MustCompile(`self\.__next_f\.push\(\[\d+,\s*"7:(.*?)"\]\)`)
	var parsed metadata.Metadata

	found := false

	doc.Find("script").EachWithBreak(func(i int, s *goquery.Selection) bool {
		matches := re.FindAllStringSubmatch(s.Text(), -1)
		log.Debug(fmt.Sprintf("Matches: %+v", matches))
		log.Debug(fmt.Sprintf("Found %d potential matches in script %d", len(matches), i+1))

		for i, m := range matches {
			if len(m) < 2 {
				log.Warn("Not enough content here. Looping though for next...")
				continue
			}

			log.Debug(fmt.Sprintf("Raw JSON attempt: %s", m[1]))

			log.Info(fmt.Sprintf("Match %d found with enough content!:", i+1))

			movieDetailsStr, ok := extractMovieDetails(m[1], log)
			if !ok {
				log.Error("Not enough movieDetails content here. Looping for next...")
				return true
			}

			log.Debug(fmt.Sprintf("Extracted movieDetails: %s", movieDetailsStr))

			err := json.Unmarshal([]byte(movieDetailsStr), &parsed)
			if err != nil {
				log.Error(fmt.Sprintf("Failed to parse movieDetails: %v", err))
				return true
			}

			// Fill nested IDs from root-level fields if present
			var raw map[string]interface{}
			if err := json.Unmarshal([]byte(movieDetailsStr), &raw); err == nil {
				if imdb, ok := raw["imdbId"].(string); ok && parsed.IDs.IMDB == "" {
					parsed.IDs.IMDB = imdb
				}
				if tmdb, ok := raw["tmdb"].(string); ok && parsed.IDs.TMDB == "" {
					parsed.IDs.TMDB = tmdb
				}
			}

			log.Success("Movie Info parsed successfully!")
			logger.LogMetadata(&parsed)

			if parsed.IDs.IMDB == "" && parsed.IDs.TMDB == "" {
				log.Warn("Neither IMDB nor TMDB ID found in script, continuing without IDs")
			}

			if parsed.Title != "" {
				log.Success("Movie details extracted from script")
				found = true
				return false
			}

		}
		return true
	})

	if !found {
		log.Error("Movie details not found in page script")
		return nil, errors.New("movie details script not found")
	}

	return &parsed, nil
}

func extractMovieDetails(raw string, log iface.Logger) (string, bool) {
	rawClean := strings.ReplaceAll(raw, `\n`, "")
	rawClean = strings.ReplaceAll(rawClean, `\"`, `"`)

	var found string

	var search func(g gjson.Result) bool
	search = func(g gjson.Result) bool {
		if g.IsObject() {
			if md := g.Get("movieDetails"); md.Exists() {
				found = md.Raw
				return false // stop search
			}
			// also check all object fields
			stop := false
			g.ForEach(func(k, v gjson.Result) bool {
				if !search(v) {
					stop = true
					return false
				}
				return true
			})
			return !stop
		} else if g.IsArray() {
			stop := false
			g.ForEach(func(_, v gjson.Result) bool {
				if !search(v) {
					stop = true
					return false
				}
				return true
			})
			return !stop
		}
		return true
	}

	search(gjson.Parse(rawClean))

	if found == "" {
		log.Error("movieDetails not found")
		return "", false
	}

	log.Debug("Extracted movieDetails JSON: " + found)
	return found, true
}
