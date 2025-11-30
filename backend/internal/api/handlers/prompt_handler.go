package handlers

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/entities"
	"chat2pay/internal/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func Prompt(cfg *yaml.Config, service service.AiService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var (
			requestBody = entities.AiAsk{}
			ctx         = c.Context()
			//log         = logger.NewLog("login_handler", cfg.Logger.Enable)
		)

		err := c.BodyParser(&requestBody)
		if err != nil {
			//log.Error(fmt.Sprintf(`error parsing request body got %s`, err))
			c.Status(http.StatusUnprocessableEntity)
			return c.JSON(presenter.ErrorResponse(err))
		}

		serv := service.Prompt(ctx, &requestBody)
		if serv.Code != 200 {
			//log.Error(fmt.Sprintf(`error on service login got %s`, serv.Errors))
			c.Status(serv.Code)
			return c.JSON(presenter.ErrorResponse(serv.Errors))
		}

		c.Status(200)
		return c.JSON(presenter.SuccessResponse(serv.Data))
	}
}
