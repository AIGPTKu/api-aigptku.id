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

	domainCt "github.com/AIGPTku/api-aigptku.id/app/controller/domain"
	"github.com/AIGPTku/api-aigptku.id/app/model"
	domainRepo "github.com/AIGPTku/api-aigptku.id/app/repository/domain"
)

func (a *repoApi) AskGPT(ctx context.Context, ask domainRepo.RequestAsk) {
	if len(ask.AskContent) == 0 {
		ask.Finish <- true
		return 
	} else if len(ask.AskContent) > 10 {
		ask.AskContent = ask.AskContent[len(ask.AskContent)-10:]
	}

	systemRequest := domainRepo.AskContent{
		Role: "system",
		Content: "You're AIGPTku Premium featured by ChatGPT Plus (4o | 4.0). use '````markdown````', don't send md if not requested by user!. use '---' opening closing if possible.",
	}

	if ask.UseDefaultSystem {
		ask.AskContent = append([]domainRepo.AskContent{
			systemRequest,
		}, ask.AskContent...)
	}

	// URL of the third-party API providing the text/event-stream
	url := "https://api.openai.com/v1/chat/completions"

	payload := map[string]any{
		"model": "gpt-4o-mini",
		"stream": true,
		"stream_options": map[string]bool{
			"include_usage": true,
		},
		"messages": ask.AskContent,
	}

	if ask.UseFunction {
		payload["functions"] = []map[string]any{
			{
				"name": "get_web_search",
				"description": "Search the web for a query",
				"parameters": map[string]any{
					"type": "object",
					"properties": map[string]any {
						"query": map[string]any{
							"type": "string",
							"description": "The search query",
						},
					},
					"required": []string{
						"query",
					},
				},
			},
			{
				"name": "image_generate",
				"description": "Generate image based on prompt",
				"parameters": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"prompt": map[string]any{
							"type": "string",
							"description": "The generate prompt. must using user language!",
						},
					},
					"required": []string{
						"prompt",
					},
				},
			},
		}
		payload["function_call"] = "auto"
	}

	body, _ := json.Marshal(payload)

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

		ask.Result <- "Maaf terjadi kesalahan. bisa ulangi lagi pertanyaannya?"
		ask.Finish <- true
		return 
	}

	// Create a scanner to read the stream
	scanner := bufio.NewScanner(resp.Body)

	assistantResponse := map[string]string{
		"role": "assistant",
		"content": "",
	}

	isFunctionCall := false
	functionCallName := ""
	functionCallArgumentsStr := ""

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
			if isFunctionCall {
				functionCallArgumentsStr += response.Choices[0].Delta.FunctionCall.Arguments
				continue
			}

			funcCall := response.Choices[0].Delta.FunctionCall.Name
			if funcCall != "" {
				isFunctionCall = true
				functionCallName = funcCall
				continue
			}

			content := response.Choices[0].Delta.Content	
			if content == "" {
				continue
			}		
			// fmt.Print(content)

			// send chan
			time.Sleep(25 * time.Millisecond)
			ask.Result <- content
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

	if isFunctionCall {
		arguments := domainCt.Arguments{}

		err = json.Unmarshal([]byte(functionCallArgumentsStr), &arguments)
		if err != nil {
			log.Println("Error unmarshalling arguments")
		}

		ask.FuncCall <- domainCt.FuncCall{
			Name: functionCallName,
			Arguments: arguments,
		}
	}

	fmt.Print("\n\n")
	ask.Finish <- true
}