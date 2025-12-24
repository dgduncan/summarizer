package summarizer

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func transcribeWhisperX(path string) {
	// 1. Configuration
	filePath := "./test_audio.mp3" // Change this to your audio file
	serverURL := "http://localhost:8080/transcribe"

	// 2. Open the file
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 3. Create Multipart Payload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		panic(err)
	}
	writer.Close()

	// 4. Create Request
	req, err := http.NewRequest("POST", serverURL, body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 5. Send Request
	client := &http.Client{Timeout: 5 * time.Minute} // Long timeout for transcription
	fmt.Println("Uploading and transcribing... (this may take a moment)")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 6. Print Response
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", string(respBody))
}
