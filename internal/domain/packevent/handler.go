package packevent

import (
	"pack-management/internal/pkg/validator"
	"time"

	"github.com/gofiber/fiber/v2"
)

type (
	handler struct {
		service Service
		app     *fiber.App
	}

	HandlerParams struct {
		App     *fiber.App `validate:"required"`
		Service Service    `validate:"required"`
	}

	CreateEventRequest struct {
		PackID      string    `json:"pack_id" validate:"required"`
		Description string    `json:"description" validate:"required"`
		Location    string    `json:"location" validate:"required"`
		Date        time.Time `json:"date" validate:"required"`
	}

	EventJSON struct {
		ID          string    `json:"id"`
		PackID      string    `json:"pack_id"`
		Description string    `json:"description"`
		Location    string    `json:"location"`
		Date        time.Time `json:"date"`
	}
)

func NewHTPPHandler(params *HandlerParams) *handler {
	params.validate()

	h := &handler{
		service: params.Service,
		app:     params.App,
	}

	group := h.app.Group("/pack_events")
	group.Post("/", h.createEvent)

	return h
}

func (p *HandlerParams) validate() {
	err := validator.ValidateStruct(p)
	if err != nil {
		panic(err)
	}
}

func (h *handler) createEvent(ctx *fiber.Ctx) error {
	payload := &CreateEventRequest{}
	if err := ctx.BodyParser(payload); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	err := validator.ValidateStruct(payload)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	_, err = h.service.CreateEvent(ctx.Context(), payload.ToEntity())
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (r *CreateEventRequest) ToEntity() *Entity {
	return &Entity{
		PackID:      r.PackID,
		Description: r.Description,
		Location:    r.Location,
		Date:        r.Date,
	}
}
