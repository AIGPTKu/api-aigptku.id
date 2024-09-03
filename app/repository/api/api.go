package api

import (
	"net/http"

	"github.com/AIGPTku/api-aigptku.id/app/repository"
)

type repoApi struct {
	client *http.Client
	gptApiKey string
	geminiApiKey string
}

func New(client *http.Client, gptApiKey, geminiApiKey string) repository.ApiInterface {
	return &repoApi{
		client,
		gptApiKey,
		geminiApiKey,
	}
}