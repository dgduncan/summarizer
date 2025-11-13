package summarizer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
)

func Summarize(ctx context.Context, transcript string) {
	u, err := url.Parse("http://ollama:11434")
	if err != nil {
		panic(err)
	}
	client := api.NewClient(u, http.DefaultClient)

	// The prompt is the key to getting good summarization
	fullPrompt := fmt.Sprintf(
		`Role: You are an expert financial analyst assistant.

Task: Analyze the following transcript from a stock picker's livestream and extract all actionable intelligence. Your output must be clear, concise, and organized for a busy investor.

Input: %s

Output Format: Provide your response using the following structure. If no information is found for a section, state "None mentioned."

Direct Actions (High Conviction)
Discussed: [List any stocks the speaker discussed.]

Buy: [List any stocks the speaker is actively buying or has a "Buy" rating on, including any reasoning.]

Sell: [List any stocks the speaker is actively selling or has a "Sell" rating on, including any reasoning.]

Long: [List any new or existing long positions mentioned.]

Short: [List any new or existing short positions mentioned.]

Position Management
Adding to: [List stocks the speaker is adding to an existing position.]

Trimming/Reducing: [List stocks the speaker is trimming or reducing.]

Closing: [List any positions the speaker has fully closed.]

Watchlist & Future Plays (Lower Conviction)
Watching: [List any stocks the speaker is adding to a watchlist or "keeping an eye on."]

Price Alerts: [List any specific price targets, entry points, or stop-losses mentioned (e.g., "Watching XYZ if it breaks $50").]`, transcript)

	req := &api.GenerateRequest{
		Model:  "llama3.1:8b", // The model you pulled locally
		Prompt: fullPrompt,
		Options: map[string]any{
			"num_ctx": 32768,
		},
		// System: `You are an assistant that takes video transcriptions and extracts the key takeaways or actions mentioned in the video. You will be told the type of video and the transcript, and you will be given a specific goals to relay in your summary`,
	}

	// Call the local Ollama API
	var result string
	if err := client.Generate(ctx, req, func(gr api.GenerateResponse) error {
		result += gr.Response
		return nil
	}); err != nil {
		panic(err)
	}

	log.Println(result)

}
