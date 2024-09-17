package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/AIGPTku/api-aigptku.id/app/model"
	"github.com/AIGPTku/api-aigptku.id/app/repository/domain"
)

func (a *repoApi) GenerateImage(ctx context.Context, content, image chan string, finish chan bool, prompt string) {
		// URL of the third-party API providing the text/event-stream
		url := "https://api.openai.com/v1/images/generations"

		body, _ := json.Marshal(map[string]any{
			"model": "dall-e-3",
			"n": 1,
			"size": "1024x1024",
			"prompt": prompt,
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
			finish <- true
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed to generate image")
			finish <- true
			return
		}
		var response model.ResponseDallE3

		err = json.Unmarshal(b, &response)
		if err != nil {
			log.Println("Failed unmarshall response dall-e 3")
			finish <- true
			return
		}

		if response.Error.Code != "" {
			log.Printf("Error system, code: %s, message: %s", response.Error.Code, response.Error.Message)
			content <- "Maaf Terjadi kesalahan, silahkan coba lagi nanti."
			finish <- true
			return
		}

		image <- response.Data[0].Url
		
		go a.AskGPT(ctx, domain.RequestAsk{
			Result: content,
			Finish: finish,
			AskContent: []domain.AskContent{
				{
					Role: "system",
					Content: fmt.Sprintf("add some opening text like 'this is ...','here are ...' or others to '%s' use same language is must!", prompt),
				},
			},
			UseDefaultSystem: false,
		})
}