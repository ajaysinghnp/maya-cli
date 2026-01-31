package metadata

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ajaysinghnp/maya-cli/utils"
)

func PromptUser(log Logger) (*Metadata, error) {
	reader := bufio.NewReader(os.Stdin)

	log.Info("Metadata required for Jellyfin formatting")

	fmt.Print("Is this a movie or series? (movie/series): ")
	t, _ := reader.ReadString('\n')
	t = strings.TrimSpace(strings.ToLower(t))

	meta := &Metadata{
		Type: MediaType(t),
	}

	fmt.Print("Title: ")
	meta.Title, _ = reader.ReadString('\n')
	meta.Title = strings.TrimSpace(meta.Title)

	fmt.Print("Year (optional): ")
	inputStr, _ := reader.ReadString('\n')
	inputStr = strings.TrimSpace(inputStr)
	meta.Year = utils.NormalizeYear(inputStr)

	fmt.Print("TMDB ID (optional): ")
	meta.IDs.TMDB, _ = reader.ReadString('\n')
	meta.IDs.TMDB = strings.TrimSpace(meta.IDs.TMDB)

	fmt.Print("IMDB ID (optional): ")
	meta.IDs.IMDB, _ = reader.ReadString('\n')
	meta.IDs.IMDB = strings.TrimSpace(meta.IDs.IMDB)

	if meta.Type == Series {
		fmt.Print("Season number: ")
		s, _ := reader.ReadString('\n')
		meta.Season, _ = strconv.Atoi(strings.TrimSpace(s))

		fmt.Print("Episode number: ")
		e, _ := reader.ReadString('\n')
		meta.Episode, _ = strconv.Atoi(strings.TrimSpace(e))
	}

	return meta, nil
}
