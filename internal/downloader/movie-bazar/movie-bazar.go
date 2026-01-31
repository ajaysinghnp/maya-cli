package moviebazar

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"strings"
	"time"

	"github.com/ajaysinghnp/maya-cli/internal/downloader/m3u8"
	"github.com/ajaysinghnp/maya-cli/internal/logger/iface"
	"github.com/ajaysinghnp/maya-cli/internal/metadata"
	"github.com/manifoldco/promptui"
)

type Options struct {
	Meta       *metadata.Metadata
	TempDir    string
	Resume     bool
	Concurrent int
	Log        iface.Logger
}

// HandleMovie handles MoviesBazar downloads with interactive source selection
func HandleMovie(opts Options) error {
	log := opts.Log

	if opts.Meta == nil || len(opts.Meta.Sources) == 0 {
		log.Error("No sources found for this movie")
		return errors.New("no sources found")
	}

	var selected metadata.Source

	prompt := promptui.Select{
		Label: "Select a source (use arrow keys or type number)",
		Items: opts.Meta.Sources,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U000027A4 {{ .Label | cyan }} ({{ .LabelTag | yellow }})",
			Inactive: "  {{ .Label }} ({{ .LabelTag | faint }})",
			Selected: "\U00002714 Selected: {{ .Label }}",
			Details: `
--------- Sources ----------
{{ "Label:" | faint }}	{{ .Label }}
{{ "LabelTag:" | faint }}	{{ .LabelTag }}
{{ "url:" | faint }}	{{ .URL }}`,
		},
		Size: 5,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}
	selected = opts.Meta.Sources[i]

	log.Success(fmt.Sprintf("Selected source: %s", selected.Label))

	// lets try to fetch the contents of the player URL
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	baseHeaders := http.Header{
		"User-Agent": []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:148.0) Gecko/20100101 Firefox/148.0",
		},
		"Referer": []string{"https://www.moviesbazar.watch/"},
		"Accept":  []string{"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
	}

	if opts.Meta.IDs.IMDB == "" && opts.Meta.IDs.TMDB == "" {
		return errors.New("IMDB/TMDB id missing, cannot resolve MoviesBazar player")
	}

	var id string

	if opts.Meta.IDs.IMDB == "" {
		id = opts.Meta.IDs.TMDB
	} else {
		id = opts.Meta.IDs.IMDB
	}

	playURL := fmt.Sprintf("https://vekna402las.com/play/%s", id)
	log.Info("Loading player page: " + playURL)

	playHTML, err := httpGet(client, playURL, baseHeaders)
	if err != nil {
		return err
	}

	fileURL, err := extractP3File(playHTML)
	if err != nil {
		return err
	}

	log.Success("Resolved playlist gateway")
	log.Debug("Playlist gateway: " + fileURL)

	// ---------- step 3: fetch playlist (.txt → m3u8) ----------
	baseHeaders.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:148.0) Gecko/20100101 Firefox/148.0")
	baseHeaders.Set("Accept", "*/*")
	baseHeaders.Set("Accept-Language", "en-US,en;q=0.9")
	baseHeaders.Set("Accept-Encoding", "gzip, deflate, br")
	baseHeaders.Set("Referer", "https://vekna402las.com/")
	baseHeaders.Set("Origin", "https://vekna402las.com")
	baseHeaders.Set("DNT", "1")
	baseHeaders.Set("Sec-Fetch-Dest", "empty")
	baseHeaders.Set("Sec-Fetch-Mode", "cors")
	baseHeaders.Set("Sec-Fetch-Site", "same-site")
	baseHeaders.Set("Connection", "keep-alive")

	playlistBody, err := httpGet(client, fileURL, baseHeaders)
	if err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Fetched playlist, length: %d", len(playlistBody)))
	log.Info(fmt.Sprintf("Body: %s", playlistBody))

	playlistURL := extractM3U8FromText(playlistBody)
	if playlistURL == "" {
		return errors.New("failed to resolve final m3u8 url")
	}

	log.Success("Final M3U8 resolved")
	log.Debug("M3U8 URL: " + playlistURL)

	// Start download using your existing m3u8 package
	return m3u8.Download(m3u8.Options{
		URL:        selected.URL,
		Output:     opts.Meta.MediaFile,
		TempDir:    opts.TempDir,
		Resume:     opts.Resume,
		Concurrent: opts.Concurrent,
		Log:        log,
	})
}

func httpGet(client *http.Client, url string, headers http.Header) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header = headers.Clone()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

func extractP3File(html string) (string, error) {
	re := regexp.MustCompile(`"file"\s*:\s*"([^"]+)"`)
	m := re.FindStringSubmatch(html)
	if len(m) < 2 {
		return "", errors.New("p3.file not found in player html")
	}

	url := strings.ReplaceAll(m[1], `\/`, `/`)
	return url, nil
}

func extractM3U8FromText(body string) string {
	// some responses are already m3u8, some redirect
	if strings.Contains(body, ".m3u8") {
		re := regexp.MustCompile(`https?://[^\s'"]+\.m3u8`)
		return re.FindString(body)
	}
	return ""
}
