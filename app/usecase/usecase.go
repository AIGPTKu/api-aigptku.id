package usecase

import (
	"context"
	"database/sql"
	"net/http"

	domainUc "github.com/AIGPTku/api-aigptku.id/app/usecase/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type InitMysql struct {
	Master *sql.DB
	Trx *sql.DB
}

type InitMongo struct {
	Master *mongo.Database
	Trx *mongo.Database
}

type InitAPI struct {
	Client *http.Client
	GPTApiKey string
	GeminiApiKey string
}

type InitRequest struct {
	Mysql InitMysql
	Mongo InitMongo
	Api InitAPI
}

type GPTUsecase interface {
	AskGPT(ctx context.Context, res chan string, finish chan bool, content []domainUc.AskContent)
}

type GeminiUsecase interface {

}