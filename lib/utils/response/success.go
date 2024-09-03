package response

import "github.com/gofiber/fiber/v2"

func Success(c *fiber.Ctx, data any, code ...int) error {
	var statusCode = 200

	if len(code) > 0 {
		statusCode = code[0]
	}

	return c.Status(statusCode).JSON(BaseResponse{
		StatusCode: statusCode,
		Data: data,
	})
}