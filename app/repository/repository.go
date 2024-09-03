package repository

import (
	"context"

	domainCt "github.com/AIGPTku/api-aigptku.id/app/controller/domain"
)

type MysqlInterface interface {
}

type MongoInterface interface {
}

type ApiInterface interface {
	AskGPT(ctx context.Context, res chan string, finish chan bool, askContent []domainCt.AskContent)
}