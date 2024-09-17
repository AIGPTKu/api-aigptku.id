package repository

import (
	"context"

	domainRepo "github.com/AIGPTku/api-aigptku.id/app/repository/domain"
)

type MysqlInterface interface {
}

type MongoInterface interface {
}

type ApiInterface interface {
	AskGPT(ctx context.Context, req domainRepo.RequestAsk)
	GenerateImage(ctx context.Context, content, image chan string, finish chan bool, prompt string)
}