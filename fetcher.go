package summarizer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	defaultDownloadDir = "./downloads"
)

func Begin(ctx context.Context, db *sql.DB, c *Config, logger *slog.Logger) error {
	for _, p := range c.Podcasts {
		logger.Debug(fmt.Sprintf("Processing the following podcasts: %s, RSS: %s", p.Name, p.RssURL))
	}

	for _, y := range c.YouTube {
		logger.Debug(fmt.Sprintf("Processing the following YouTube channels: %s, RSS: %s", y.Name, y.RssURL))
	}

	// process podcasts
	if err := processPodcasts(ctx, db, c.Podcasts, logger); err != nil {
		return fmt.Errorf("failed to process podcasts: %w", err)
	}

	// process YouTube channels
	if err := processYouTubeChannels(ctx, db, c.YouTube, logger); err != nil {
		return fmt.Errorf("failed to process YouTube channels: %w", err)
	}

	return nil
}

func processPodcasts(ctx context.Context, db *sql.DB, podcasts []PodcastConfig, logger *slog.Logger) error {
	for _, p := range podcasts {
		logger.Debug(fmt.Sprintf("fetching rss feed %s", p.RssURL))
		feed := GetRSSFeed(p.RssURL)
		latest := GetLatestEpisode(feed)
		fmt.Println(latest.Title)

		// check if podcast already downloaded
		logger.Debug("checking if download already exists")
		mp3Path := filepath.Join(defaultDownloadDir, fmt.Sprintf(defaultFileType, latest.GUID))
		if _, err := os.Stat(mp3Path); err == nil {
			logger.Debug(fmt.Sprintf("podcast %s already downloaded, not downloading...", latest.Title))
		} else {
			// download podcast
			logger.Info(fmt.Sprintf("downloading podcast: %s", latest.Title))
			mp3Path, err = DownloadPodcast(ctx, latest.Enclosures[0].URL, latest.GUID, defaultDownloadDir)
			if err != nil {
				return fmt.Errorf("failed to download podcast %s: %w", latest.Title, err)
			}
		}
		logger.Info(mp3Path)

		wavPath, err := ConvertToWav(mp3Path, filepath.Join(defaultDownloadDir, fmt.Sprintf("%s.wav", latest.GUID)))
		if err != nil {
			return fmt.Errorf("failed to convert podcast %s to wav: %w", latest.Title, err)
		}

		logger.Info(fmt.Sprintf("successfully downloaded podcast: %s", latest.Title))

		transcription := GenerateTranscription(wavPath)
		log.Println(transcription)
	}

	return nil

}

func processYouTubeChannels(ctx context.Context, db *sql.DB, channels []YouTubeConfig, logger *slog.Logger) error {
	return nil
}
