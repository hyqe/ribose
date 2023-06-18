package fit

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/hyqe/ribose/internal/fit/codes"
	"github.com/hyqe/ribose/internal/fit/status"
)

// Fiber builds a POST/JSON generic handler
func Fiber[I, O any](fn func(ctx context.Context, in I) (O, status.Status)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var in I
		err := c.BodyParser(&in)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		switch out, status := fn(c.Context(), in); status.Code {
		case codes.OK, codes.Created, codes.Accepted:
			return c.JSON(out)
		default:
			return fiber.NewError(int(status.Code), status.Error())
		}
	}
}
