package summarizer

import "fmt"

const ()

func GenerateTranscription(path string) string {
	resp, err := execute(commandWhisper, []string{"-m", "whisper.cpp/models/ggml-base.en.bin", "-f", path, "-otxt", "-of", "transcription", "--no-timestamps"})
	if err != nil {
		return fmt.Sprintf("Error generating version: %s : %v", string(resp), err)
	}

	return fmt.Sprintf("Transcription: %s", resp)
}
