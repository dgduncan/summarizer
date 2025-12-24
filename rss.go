package summarizer

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	defaultFileType = "%s.mp3"
)

const (
	profGRSSFeed = "https://feeds.megaphone.fm/profgmarkets"
)

func GetRSSFeed(url string) *gofeed.Feed {

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(profGRSSFeed)

	return feed
}

func GetLatestEpisode(feed *gofeed.Feed) *gofeed.Item {
	if len(feed.Items) == 0 {
		return nil
	}
	return feed.Items[0]
}

func DownloadPodcast(ctx context.Context, url, name, dir string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Overcast/2025.1 (macOS) GitHubCopilot")
	req.Header.Set("Accept", "audio/mpeg,*/*;q=0.8")

	client := &http.Client{
		Timeout: 2 * time.Minute,
		// default CheckRedirect follows redirects; keep it
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	outPath := filepath.Join(dir, fmt.Sprintf(defaultFileType, name))
	f, err := os.Create(outPath)
	if err != nil {
		log.Println("Error creating file:", err)
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}
	return outPath, nil
}
