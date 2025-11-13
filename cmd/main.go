package main

import (
	"context"
	"log"
	"os"

	"github.com/dgduncan/summarizer"
)

func main() {

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

	data, err := os.ReadFile("transcriptions/joseph_carlson.txt")
	if err != nil {
		log.Fatalf("Failed to read transcription.txt: %v", err)
	}

	log.Println(string(data))

	summarizer.Summarize(context.Background(), string(data))

}
