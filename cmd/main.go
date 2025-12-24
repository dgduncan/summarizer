package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/dgduncan/summarizer"
)

const (
	defaulEnvConfig = "CONFIG_PATH"
)

func run(ctx context.Context, w io.Writer, args []string) error {
	// handle cmd flags
	cpath := summarizer.CmdFlags(args)

	if cpath == nil || *cpath == "" {
		if cepath := os.Getenv(defaulEnvConfig); cepath == "" {
			return fmt.Errorf("config path not provided via -c flag or %s env var", defaulEnvConfig)
		} else {
			cpath = &cepath
		}
	}

	c, err := summarizer.LoadConfig(*cpath)
	if err != nil {
		return errors.Join(errors.New("Failed to load config"), err)
	}

	logger := c.Logging.ToLogger(w)

	logger.Debug("starting summarizer...")

	// init db
	db, err := summarizer.NewSQLiteDB()
	if err != nil {
		return errors.Join(errors.New("failed to initialize SQlite DB"), err)
	}
	defer db.Close()

	if err := summarizer.Begin(ctx, db, c, logger); err != nil {
		return errors.New("failed to run summarizer")
	}

	// v, _ := summarizer.Get(summarizer.SummarysBucket, "7e583c7e-1fbc-11f0-af56-cb7dc28e9c4c", db)
	// fmt.Println(v)

	// iterate through podcasts first
	// for _, p := range c.Podcasts {
	// 	log.Printf("Processing podcast: %s", p.Name)

	// 	feed := summarizer.GetRSSFeed(p.RssURL)
	// 	// feed.Items = feed.Items[:min(1, len(feed.Items))] // limit to latest 5 episodes

	// 	// check which episodes have not been process yet
	// 	unprocessed := make([]*gofeed.Item, 0)
	// 	for _, it := range feed.Items {
	// 		key := it.GUID
	// 		log.Println(it.Title, key)

	// 		// val, err := summarizer.Get(summarizer.ShowsBucket, key, db)
	// 		// if err != nil {
	// 		// 	panic(err)
	// 		// }
	// 		// if val != "" {
	// 		// 	log.Println("Found existing entry for key:", key, "skipping...")
	// 		// 	continue
	// 		// }
	// 		unprocessed = append(unprocessed, it)
	// 	}
	// 	log.Printf("Found %d unprocessed episodes for podcast: %s", len(unprocessed), p.Name)

	// 	// process unprocessed episodes
	// 	for _, it := range unprocessed {
	// 		log.Printf("Downloading episode: %s", it.Title)
	// 		if err := summarizer.DownloadMP3(ctx, it.Enclosures[0].URL, it.GUID); err != nil {
	// 			log.Fatalf("Failed to download episode %s: %v", it.Title, err)
	// 		}
	// 		log.Printf("Creating transcription for episode: %s", it.Title)
	// 		transcription := summarizer.GenerateTranscription(fmt.Sprintf("%s.mp3", it.GUID))

	// 		fmt.Println(transcription)

	// 		// if err := summarizer.Set(summarizer.SummarysBucket, it.GUID, transcription, db); err != nil {
	// 		// 	log.Fatalf("Failed to store transcription for episode %s: %v", it.Title, err)
	// 		// }

	// 		log.Printf("Finished processing episode: %s", it.Title)

	// 		// 	if err := summarizer.Set(summarizer.ShowsBucket, it.GUID, "true", db); err != nil {
	// 		// 		log.Fatalf("Failed to mark episode %s as processed: %v", it.Title, err)
	// 		// 	}
	// 	}
	// }

	// log.Println("beginning download...")
	// resp, err := summarizer.DownloadAudio("youtube.com/watch?v=z4U9Ih4f-PU")
	// if err != nil {
	// 	log.Println(err, resp)
	// 	return
	// }
	// log.Println("finished download...")
	// // log.Println("Download Audio Response:", resp)

	// log.Println("creating transcription...")
	// fmt.Println(summarizer.GenerateTranscription("my_audio.wav"))
	// log.Println("finished transcription...")

	// log.Println(resp2)

	// data, err := os.ReadFile("transcriptions/joseph_carlson.txt")
	// if err != nil {
	// 	log.Fatalf("Failed to read transcription.txt: %v", err)
	// }

	// log.Println(string(data))

	// summarizer.Summarize(context.Background(), string(data))

	// summarizer.DownloadMP3(it.Links[0])
	// for _, it := range feed.Items {
	// 	key := it.GUID // human-readable
	// 	log.Println("episode key:", key)
	// 	log.Println(it.Published)
	// 	log.Println(it.PublishedParsed)

	// 	val, err := summarizer.Get(key, db)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if val != "" {
	// 		log.Println("Found existing entry for key:", key, "skipping...")
	// 		continue
	// 	}

	// 	// download and process the episode...
	// 	log.Println("Storing entry for key:", key)
	// 	if err = summarizer.Set(key, "done", db); err != nil {
	// 		panic(err)
	// 	}
	// }
	return nil
}

func main() {
	ctx := context.Background()

	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
