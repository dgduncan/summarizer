package summarizer

import (
	"fmt"
	"runtime"
	"time"
)

func GenerateTranscription(path string) string {
	now := time.Now()
	fmt.Printf("Generating transcription at %s for file %s\n", now.Format(time.RFC3339), path)
	resp, err := execute("ffmpeg", []string{"-i", path, "-ar", "16000", "-ac", "1", "output.wav"})
	if err != nil {
		return fmt.Sprintf("Error converting to wav: %s : %v", string(resp), err)
	}
	fmt.Printf("Conversion to wav completed at %s\n", time.Now().Format(time.RFC3339))

	now = time.Now()
	fmt.Printf("Starting transcription at %s\n", now.Format(time.RFC3339))
	resp, err = execute(commandWhisper, []string{"-m", "whisper.cpp/models/ggml-small.en-q8_0.bin", "-f", "output.wav", "-otxt", "-of", "transcription", "--no-timestamps", "-t", fmt.Sprint(runtime.NumCPU())})
	if err != nil {
		return fmt.Sprintf("Error generating version: %s : %v", string(resp), err)
	}
	fmt.Printf("Transcription completed at %s\n", time.Since(now))

	return fmt.Sprintf("Transcription: %s", resp)
}
