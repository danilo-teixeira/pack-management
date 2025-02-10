package pack

import (
	"pack-management/internal/pkg/validator"
	"pack-management/internal/services/person"
	"time"

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
		Description           string    `json:"description" validate:"required"`
		Status                Status    `json:"status" validate:"required"`
		ReceiverName          string    `json:"recipient" validate:"required"`
		SenderName            string    `json:"sender" validate:"required"`
		EstimatedDeliveryDate time.Time `json:"estimated_delivery_date" validate:"required"`
	}

	CreatePackResponse struct {
		ID           string `json:"id"`
		Description  string `json:"description"`
		Status       Status `json:"status"`
		ReceiverName string `json:"recipient"`
		SenderName   string `json:"sender"`
		CreatedAt    string `json:"created_at"`
		UpdateAt     string `json:"updated_at"`
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
		return err
	}

	err := validator.ValidateStruct(payload)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest) // TODO: implement error handler
	}

	pack, err := h.service.CreatePack(ctx.Context(), payload.ToEntity())
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError) // TODO: implement error handler
	}

	response := &CreatePackResponse{}

	return ctx.Status(fiber.StatusCreated).JSON(response.FromEntity(pack))
}

func (r *CreatePackRequest) ToEntity() *Entity {
	return &Entity{
		Description:           r.Description,
		Status:                r.Status,
		EstimatedDeliveryDate: r.EstimatedDeliveryDate,
		Receiver: &person.Entity{
			Name: r.ReceiverName,
		},
		Sender: &person.Entity{
			Name: r.SenderName,
		},
	}
}

func (r *CreatePackResponse) FromEntity(pack *Entity) *CreatePackResponse {
	return &CreatePackResponse{
		ID:           pack.ID,
		Description:  pack.Description,
		Status:       pack.Status,
		ReceiverName: pack.Receiver.Name,
		SenderName:   pack.Sender.Name,
		CreatedAt:    pack.CreatedAt.String(),
		UpdateAt:     pack.UpdatedAt.String(),
	}
}
