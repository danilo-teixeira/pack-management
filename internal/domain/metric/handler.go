package metric

import (
	"pack-management/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	HandlerParams struct {
		App *fiber.App `validate:"required"`
	}

	handler struct {
		app *fiber.App `validate:"required"`
	}
)

func NewHTPPHandler(params *HandlerParams) *handler {
	params.validate()

	h := &handler{
		app: params.App,
	}

	h.app.Get("/metrics", h.listMetrics)

	return h
}

func (p *HandlerParams) validate() {
	err := validator.ValidateStruct(p)
	if err != nil {
		panic(err)
	}
}

func (h *handler) listMetrics(ctx *fiber.Ctx) error {
	return adaptor.HTTPHandler(promhttp.Handler())(ctx)
}
