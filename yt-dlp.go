package summarizer

import (
	"fmt"
	"os/exec"
)

const (
	commandYTDLP   = "yt-dlp"
	commandWhisper = "./whisper.cpp/build/bin/whisper-cli" // Adjust path as necessary
)

func execute(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("error executing command: %v", err)
	}

	return string(output), nil
}

func CheckVersion() string {
	resp, err := execute(commandYTDLP, []string{"--version"})
	if err != nil {
		return fmt.Sprintf("Error checking yt-dlp version: %v", err)
	}
	return fmt.Sprintf("yt-dlp version: %s", resp)
}

func DownloadAudio(url string) (string, error) {
	resp, err := execute(commandYTDLP, []string{"-U", "-x", "--audio-format", "wav", "-o", "my_audio.%(ext)s", url})
	if err != nil {
		return resp, fmt.Errorf("error downloading audio: %v", err)
	}
	return resp, nil
}
