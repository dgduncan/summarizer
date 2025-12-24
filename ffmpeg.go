package summarizer

import (
	"errors"
	"fmt"
	"time"
)

// ExtractAudio uses ffmpeg to convert the input audio file to a WAV file with 16kHz sample rate and mono channel.
// PLACEHOLDER ONLY
func ConvertToWav(i, o string) (string, error) {
	now := time.Now()
	_, err := execute("ffmpeg", []string{"-y", "-i", i, "-ar", "16000", "-ac", "1", o})
	if err != nil {
		return "", errors.Join(errors.New("Error converting to wav"), err)
	}
	fmt.Printf("Conversion to wav completed after %s\n", time.Since(now))

	return o, nil
}
