package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dgduncan/summarizer"
)

func main() {
	ctx := context.Background()

	// log.Println("starting summarizer...")

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

	db, err := summarizer.Open()
	if err != nil {
		panic(err)
	}
	defer summarizer.Close(db)

	feed := summarizer.GetRSSFeed()
	it := summarizer.GetLatestEpisode(feed)

	log.Println("Latest episode title:", it.Title)

	log.Println("downloading latest episode...", it.Enclosures[0].URL)
	if err := summarizer.DownloadMP3(ctx, it.Enclosures[0].URL, "latest_episode"); err != nil {
		log.Fatalf("Failed to download latest episode: %v", err)
	}
	log.Println("finished downloading latest episode...")

	log.Println("creating transcription...")
	fmt.Println(summarizer.GenerateTranscription("latest_episode.mp3"))
	log.Println("finished transcription...")

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

}
