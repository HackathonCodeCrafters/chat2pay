package handlers

import (
	"chat2pay/config/yaml"
	"chat2pay/internal/api/presenter"
	"chat2pay/internal/consts"
	"chat2pay/internal/helper"
	"chat2pay/internal/pkg/logger"
	"chat2pay/internal/service"
	"github.com/gofiber/fiber/v2"
)

func GetProfile(cfg *yaml.Config, service service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var (
			ctx = c.Context()
			log = logger.NewLog("get_profile_handler", cfg.Logger.Enable)
		)

		//log.Info(fmt.Sprintf(`start service login for user %s`, requestBody.Email))
		serv := service.GetProfile(ctx, helper.InterfaceToString(c.Locals(consts.UserId)))
		if serv.Code != 200 {
			//log.Error(fmt.Sprintf(`error on service login got %s`, serv.Errors))
			c.Status(serv.Code)
			return c.JSON(presenter.ErrorResponse(serv.Errors))
		}

		log.Info("get profile success")

		c.Status(200)
		return c.JSON(presenter.SuccessResponse(serv.Data))
	}
}
