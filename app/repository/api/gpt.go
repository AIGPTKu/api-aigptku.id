package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AIGPTku/api-aigptku.id/app/model"
	domainRepo "github.com/AIGPTku/api-aigptku.id/app/repository/domain"
)

func (a *repoApi) AskGPT(ctx context.Context, res chan string, finish chan bool, askContent []domainRepo.AskContent) {
	if len(askContent) == 0 {
		finish <- true
		return 
	} else if len(askContent) > 3 {
		askContent = askContent[len(askContent)-3:]
	}

	systemRequest := domainRepo.AskContent{
		Role: "system",
		Content: "You're AIGPTku Premium featured by ChatGPT Plus (4o | 4.0). answer more detail for copywriting and answer briefly for programming. use '````markdown````', don't send md if not requested by user!. use '---' opening closing if possible.",
	}

	askContent = append([]domainRepo.AskContent{
		systemRequest,
	}, askContent...)

	// URL of the third-party API providing the text/event-stream
	url := "https://api.openai.com/v1/chat/completions"

	body, _ := json.Marshal(map[string]any{
		"model": "gpt-4o-mini",
		"stream": true,
		"stream_options": map[string]bool{
			"include_usage": true,
		},
		"messages": askContent,
	})

	// Make a GET request to the third-party API
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		log.Println("Error creating request")
		return 
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + a.gptApiKey) // Replace with actual token

	// Send the request
	resp, err := a.client.Do(req)

	if err != nil {
		log.Println("Failed to connect to the stream")
		return
	}
	defer resp.Body.Close()

	// Check if the response is a text/event-stream
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/event-stream") {
		log.Println(resp.Header.Get("Content-Type"))
		log.Println("Invalid content type")
		body, _ := io.ReadAll(resp.Body)
		log.Println(string(body))

		res <- "Maaf terjadi kesalahan. bisa ulangi lagi pertanyaannya?"
		finish <- true
		return 
	}

	// Create a scanner to read the stream
	scanner := bufio.NewScanner(resp.Body)

	assistantResponse := map[string]string{
		"role": "assistant",
		"content": "",
	}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, "data: ", "", 1)

		// Skip empty lines
		if line == "" || line == "[DONE]" {
			continue
		}

		// Parse the JSON object
		var response model.ResponseGPT
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			log.Println("Error parsing JSON:", err)
			log.Println(strings.Replace(line, "data: ", "", 1))
			continue
		}

		// Extract and process the content
		if len(response.Choices) > 0 {
			content := response.Choices[0].Delta.Content	
			if content == "" {
				continue
			}		
			// fmt.Print(content)

			// send chan
			time.Sleep(25 * time.Millisecond)
			res <- content
			assistantResponse["content"] += content
		}

		if response.Usage.TotalTokens != 0 {
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("\nPromt: %d Completion: %d Total Tokens used: %d\n", response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading stream:", err)
		return
	}

	fmt.Print("\n\n")
	finish <- true
}