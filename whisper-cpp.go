package summarizer

import (
	"fmt"
	"runtime"
	"time"
)

func GenerateTranscription(path string) string {
	now := time.Now()
	fmt.Printf("Starting transcription at %s\n", now.Format(time.RFC3339))
	resp, err := execute(commandWhisper, []string{"-m", "whisper.cpp/models/ggml-small.en-tdrz.bin", "-f", path, "-otxt", "-of", "transcription", "--no-timestamps", "-t", fmt.Sprint(runtime.NumCPU())})
	if err != nil {
		return fmt.Sprintf("Error generating version: %s : %v", string(resp), err)
	}
	fmt.Printf("Transcription completed at %s\n", time.Since(now))

	return fmt.Sprintf("Transcription: %s", resp)
}
