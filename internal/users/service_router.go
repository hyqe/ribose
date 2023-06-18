package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hyqe/ribose/internal/fit"
)

func (s *Service) Router() *fiber.App {
	app := fiber.New()
	app.Post("/Create", fit.Fiber(s.Create))
	app.Post("/UpdateByUUID", fit.Fiber(s.UpdateByUUID))
	app.Post("/GetByUUID", fit.Fiber(s.GetByUUID))
	app.Post("/GetByEmail", fit.Fiber(s.GetByEmail))
	app.Post("/DeleteByUUID", fit.Fiber(s.DeleteByUUID))
	return app
}
