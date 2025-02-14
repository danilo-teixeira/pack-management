package pack

import (
	"errors"
	"pack-management/internal/domain/packevent"
	"pack-management/internal/domain/person"
	"pack-management/internal/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type (
	Handler interface{}

	handler struct {
		service Service
		app     *fiber.App
	}

	HandlerParams struct {
		App     *fiber.App `validate:"required"`
		Service Service    `validate:"required"`
	}

	CreatePackRequest struct {
		Description           string `json:"description" validate:"required"`
		ReceiverName          string `json:"recipient" validate:"required"`
		SenderName            string `json:"sender" validate:"required"`
		EstimatedDeliveryDate string `json:"estimated_delivery_date" validate:"required,datetime=2006-01-02"`
	}

	PackIDParam struct {
		ID string `params:"id"`
	}

	UpdatePackStatusRequest struct {
		PackIDParam
		Status Status `json:"status" validate:"required,oneof=CREATED IN_TRANSIT DELIVERED"`
	}

	PackJSON struct {
		ID           string                `json:"id"`
		Description  string                `json:"description"`
		Status       Status                `json:"status"`
		ReceiverName string                `json:"recipient"`
		SenderName   string                `json:"sender"`
		CreatedAt    string                `json:"created_at"`
		UpdateAt     string                `json:"updated_at"`
		DeliveredAt  string                `json:"delivered_at,omitempty"`
		CanceledAt   string                `json:"canceled_at,omitempty"`
		Events       []packevent.EventJSON `json:"events,omitempty"`
	}
)

func NewHTPPHandler(params *HandlerParams) Handler {
	params.validate()

	h := &handler{
		service: params.Service,
		app:     params.App,
	}

	group := h.app.Group("/packs")
	group.Post("/", h.createPack)
	group.Get("/:id", h.getPackByID)
	group.Patch("/:id", h.updatePackStatusByID)
	group.Post("/:id/cancel", h.cancelPackStatusByID)

	return h
}

func (p *HandlerParams) validate() {
	err := validator.ValidateStruct(p)
	if err != nil {
		panic(err)
	}
}

func (h *handler) createPack(ctx *fiber.Ctx) error {
	payload := &CreatePackRequest{}
	if err := ctx.BodyParser(payload); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	err := validator.ValidateStruct(payload)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	pack, err := h.service.CreatePack(ctx.Context(), payload.ToEntity())
	if err != nil {
		return h.errorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(h.packEntityToJSON(pack))
}

func (h handler) getPackByID(ctx *fiber.Ctx) error {
	params := &PackIDParam{}
	if err := ctx.ParamsParser(params); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	pack, err := h.service.GetPackByID(ctx.Context(), params.ID)
	if err != nil {
		return h.errorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(h.packEntityToJSON(pack))
}

func (h *handler) updatePackStatusByID(ctx *fiber.Ctx) error {
	params := &UpdatePackStatusRequest{}
	if err := ctx.ParamsParser(params); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	payload := &UpdatePackStatusRequest{}
	if err := ctx.BodyParser(payload); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	err := validator.ValidateStruct(payload)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	pack, err := h.service.UpdatePackStatusByID(ctx.Context(), params.ID, payload.ToEntity())
	if err != nil {
		return h.errorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(h.packEntityToJSON(pack))
}

func (h *handler) cancelPackStatusByID(ctx *fiber.Ctx) error {
	params := &PackIDParam{}
	if err := ctx.ParamsParser(params); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	pack, err := h.service.CancelPackStatusByID(ctx.Context(), params.ID)
	if err != nil {
		return h.errorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(h.packEntityToJSON(pack))
}

func (h *handler) errorHandler(ctx *fiber.Ctx, err error) error {
	if errors.Is(err, ErrPackNotFound) {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	if errors.Is(err, ErrStatusInvalid) ||
		errors.Is(err, ErrCannotCancel) {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.SendStatus(fiber.StatusInternalServerError)
}

func (r *UpdatePackStatusRequest) ToEntity() *Entity {
	return &Entity{
		Status: r.Status,
	}
}

func (r *CreatePackRequest) ToEntity() *Entity {
	return &Entity{
		Description:           r.Description,
		EstimatedDeliveryDate: r.EstimatedDeliveryDate,
		Receiver: &person.Entity{
			Name: r.ReceiverName,
		},
		Sender: &person.Entity{
			Name: r.SenderName,
		},
	}
}

func (h *handler) packEntityToJSON(pack *Entity) *PackJSON {
	if pack == nil {
		return nil
	}

	resp := &PackJSON{
		ID:           pack.ID,
		Description:  pack.Description,
		Status:       pack.Status,
		ReceiverName: pack.Receiver.Name,
		SenderName:   pack.Sender.Name,
		CreatedAt:    pack.CreatedAt.String(),
		UpdateAt:     pack.UpdatedAt.String(),
	}

	if pack.DeliveredAt != nil {
		resp.DeliveredAt = pack.DeliveredAt.String()
	}

	if pack.CanceledAt != nil {
		resp.CanceledAt = pack.CanceledAt.String()
	}

	if len(pack.Events) > 0 {
		for _, event := range pack.Events {
			resp.Events = append(resp.Events, packevent.EventJSON{
				ID:          event.ID,
				PackID:      event.PackID,
				Description: event.Description,
				Date:        event.Date.String(),
				CreatedAt:   event.CreatedAt.String(),
			})
		}
	}

	return resp
}
