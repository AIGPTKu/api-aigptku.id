package response

import "github.com/gofiber/fiber/v2"

func Error(c *fiber.Ctx, code int, message string, data ...any) error {
	var optData any
	if len(data) > 0 {
		optData = data[0]
	}
	return c.Status(code).JSON(BaseResponse{
		StatusCode: code,
		ErrorMessage: message,
		Data: optData,
	})
}