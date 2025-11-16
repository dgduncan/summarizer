package summarizer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mmcdole/gofeed"
)

const (
	defaultOutPath = "%s.mp3"
)

const (
	profGRSSFeed = "https://feeds.megaphone.fm/profgmarkets"
)

func GetRSSFeed() *gofeed.Feed {

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(profGRSSFeed)

	b, _ := json.Marshal(feed)

	os.WriteFile("dump.json", b, 0644)

	return feed
}

func GetLatestEpisode(feed *gofeed.Feed) *gofeed.Item {
	if len(feed.Items) == 0 {
		return nil
	}
	return feed.Items[0]
}

func DownloadMP3(ctx context.Context, url, name string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Overcast/2025.1 (macOS) GitHubCopilot")
	req.Header.Set("Accept", "audio/mpeg,*/*;q=0.8")

	client := &http.Client{
		Timeout: 2 * time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Println("redirected to:", req.URL.String())
			return nil
		},
		// default CheckRedirect follows redirects; keep it
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	f, err := os.Create(fmt.Sprintf(defaultOutPath, name))
	if err != nil {
		log.Println("Error creating file:", err)
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
