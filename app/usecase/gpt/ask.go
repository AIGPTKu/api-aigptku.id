package gpt

import (
	"context"

	domainUc "github.com/AIGPTku/api-aigptku.id/app/usecase/domain"
)

func (u *gptUsecase) AskGPT(ctx context.Context, res chan string, finish chan bool, content []domainUc.AskContent) {
	u.api.AskGPT(ctx, res, finish, content)
}