package gpt

import (
	"context"

	domainCt "github.com/AIGPTku/api-aigptku.id/app/controller/domain"
	domainRepo "github.com/AIGPTku/api-aigptku.id/app/repository/domain"
)

func (u *gptUsecase) HandleFunctionText(ctx context.Context, content chan string, finish chan bool, f domainCt.FuncCall) {
	switch f.Name {
	case "about_me":
		{
			u.api.AskGPT(ctx, domainRepo.RequestAsk{
				FuncCall: nil,
				Result: content,
				Finish: finish,
				AskContent: []domainCt.AskContent{
					{
						Role: "system",
						Content: "Use one of this information: ['this app is AIGPTku Premium featured with ChatGPT Plus 4o or 4.0', 'for subscription currently not available and will be available soon', 'the price is Rp49.000 for one month or Rp49.000 for 1 million tokens or Rp79.000 for 100 image credits or Rp99.000 for 200 image credits', 'feature unlimited chat based on subs packet, but not spam']",
					},
					{
						Role: "user",
						Content: f.Arguments.Ask,
					},
				},
				UseDefaultSystem: false,
				UseFunction: false,
			})
		}
	}
}