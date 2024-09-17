package gpt

import "context"

func (u *gptUsecase) GenerateImage(ctx context.Context, content, image chan string, finish chan bool, prompt string) {
	u.api.GenerateImage(ctx, content, image, finish, prompt)
}