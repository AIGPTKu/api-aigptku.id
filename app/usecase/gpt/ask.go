package gpt

import (
	"context"

	domainRepo "github.com/AIGPTku/api-aigptku.id/app/repository/domain"
)

func (u *gptUsecase) AskGPT(ctx context.Context, ask domainRepo.RequestAsk) {
	u.api.AskGPT(ctx, domainRepo.RequestAsk{
		FuncCall: ask.FuncCall,
		Result: ask.Result,
        Finish: ask.Finish,
		Abort: ask.Abort,
        AskContent: ask.AskContent,
		UseDefaultSystem: true,
		UseFunction: true,
	})
}