package pack

import (
	"pack-management/internal/domain/person"
	"pack-management/internal/pkg/cerrors"
	"pack-management/internal/pkg/pagination"
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

	CreatePackRequest struct {
		Description           string `json:"description" validate:"required"`
		ReceiverName          string `json:"recipient" validate:"required"`
		SenderName            string `json:"sender" validate:"required"`
		EstimatedDeliveryDate string `json:"estimated_delivery_date" validate:"required,datetime=2006-01-02"`
	}

	GetPackQuery struct {
		WithEvents bool `query:"with_events"`
	}

	PackIDParam struct {
		ID string `params:"id"`
	}

	ListPackQuery struct {
		SenderName   *string `query:"sender_name"`
		ReceiverName *string `query:"recipient_name"`
		PageSize     int     `query:"page_size"`
		PageCursor   *string `query:"page_cursor"`
	}

	ListPackJSON struct {
		Items    []*PackJSON         `json:"items"`
		Metadata pagination.Metadata `json:"metadata"`
	}

	UpdatePackStatusRequest struct {
		PackIDParam
		Status Status `json:"status" validate:"required,oneof=CREATED IN_TRANSIT DELIVERED"`
	}

	PackJSON struct {
		ID           string      `json:"id"`
		Description  string      `json:"description"`
		Status       Status      `json:"status"`
		ReceiverName string      `json:"recipient"`
		SenderName   string      `json:"sender"`
		CreatedAt    time.Time   `json:"created_at"`
		UpdateAt     time.Time   `json:"updated_at"`
		DeliveredAt  *time.Time  `json:"delivered_at,omitempty"`
		CanceledAt   *time.Time  `json:"canceled_at,omitempty"`
		Events       []EventJSON `json:"events,omitempty"`
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

	group := h.app.Group("/packs")
	group.Post("/", h.createPack)
	group.Get("/", h.listPacks)
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

func (h *handler) listPacks(ctx *fiber.Ctx) error {
	queries := &ListPackQuery{}
	if err := ctx.QueryParser(queries); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	filters := &ListFilters{
		SenderName:   queries.SenderName,
		ReceiverName: queries.ReceiverName,
		PageSize:     queries.PageSize,
		PageCursor:   queries.PageCursor,
	}

	packs, metadata, err := h.service.ListPacks(ctx.Context(), filters)
	if err != nil {
		return h.errorHandler(ctx, err)
	}

	packsJSON := make([]*PackJSON, 0, len(packs))
	for _, pack := range packs {
		packsJSON = append(packsJSON, h.packEntityToJSON(pack))
	}

	resp := &ListPackJSON{
		Items: packsJSON,
		Metadata: pagination.Metadata{
			PageSize:   metadata.PageSize,
			NextCursor: metadata.NextCursor,
			PrevCursor: metadata.PrevCursor,
		},
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

func (h handler) getPackByID(ctx *fiber.Ctx) error {
	params := &PackIDParam{}
	if err := ctx.ParamsParser(params); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	queries := &GetPackQuery{}
	if err := ctx.QueryParser(queries); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	pack, err := h.service.GetPackByID(ctx.Context(), params.ID, queries.WithEvents)
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
	if cerrors.Is(err, ErrPackNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(err)
	}

	if cerrors.Is(err, ErrStatusInvalid) ||
		cerrors.Is(err, ErrCannotCancel) {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
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
		CreatedAt:    pack.CreatedAt,
		UpdateAt:     pack.UpdatedAt,
	}

	if pack.DeliveredAt != nil {
		resp.DeliveredAt = pack.DeliveredAt
	}

	if pack.CanceledAt != nil {
		resp.CanceledAt = pack.CanceledAt
	}

	if len(pack.Events) > 0 {
		for _, event := range pack.Events {
			resp.Events = append(resp.Events, EventJSON{
				ID:          event.ID,
				PackID:      event.PackID,
				Description: event.Description,
				Location:    event.Location,
				Date:        event.Date,
			})
		}
	}

	return resp
}
