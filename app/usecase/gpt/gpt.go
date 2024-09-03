package gpt

import (
	"github.com/AIGPTku/api-aigptku.id/app/repository"
	"github.com/AIGPTku/api-aigptku.id/app/repository/api"
	"github.com/AIGPTku/api-aigptku.id/app/repository/mysql"
	"github.com/AIGPTku/api-aigptku.id/app/usecase"
)

type gptUsecase struct {
	mysql repository.MysqlInterface
	api repository.ApiInterface
}

func NewUsecase(r usecase.InitRequest) usecase.GPTUsecase {
	return &gptUsecase{
		mysql: mysql.New(r.Mysql.Master, r.Mysql.Trx),
		api: api.New(r.Api.Client, r.Api.GPTApiKey, r.Api.GeminiApiKey),
	}
}