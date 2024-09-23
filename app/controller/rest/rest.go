package rest

import (
	"github.com/AIGPTku/api-aigptku.id/app/usecase"
	"github.com/AIGPTku/api-aigptku.id/app/usecase/gpt"
	"github.com/gofiber/fiber/v2"
)

type restHandler struct {
	app *fiber.App
	uc *uc
}

type uc struct {
	gpt usecase.GPTUsecase
	// gemini usecase.GeminiUsecase
}

type InitRestHandler struct {
	usecase.InitRequest
}

func NewRestHandler(app *fiber.App, r *InitRestHandler) *restHandler {
	return &restHandler{
		app: app,
		uc: &uc{
            gpt: gpt.NewUsecase(r.InitRequest),
            // gemini:,
        },
	}
}

func (r *restHandler) RegisterRoute() {
	v1 := r.app.Group("/v1")

	v1.Post("/generative", r.ask)
	v1.Post("/generative/image", r.generateImage)
	v1.Post("/upload/temp", r.uploadFileTemp)
}