package deeplink_handler

import (
	"deeplink-bff/bff/internal/adapters/handler/dto"
	"deeplink-bff/bff/internal/core/ports"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	deeplinkService ports.DeeplinkService
}

func NewHandler(deeplinkService ports.DeeplinkService) *Handler {
	return &Handler{
		deeplinkService,
	}
}

// @Summary	get deeplink List
// @Schemes
// @Description	endpoint for get deeplink list
// @Tags			deeplink
// @Accept			application/json
// @Produce		json
// @Success		200	{object}	dto.GetDeeplinkListResponse
// @Router			/v1/deeplink [get]
// @Security		Authorization
func (h *Handler) GetDeeplinkList(c *fiber.Ctx) error {
	ctx := c.UserContext()

	slog.InfoContext(ctx, "Calling GetDeeplinkList in handler", slog.Any("test1", "testinfo1"))

	deeplinks, err := h.deeplinkService.GetDeeplinkList(ctx)
	if err != nil {
		// TODO: using api error standard response
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// TODO: using api standard response
	return c.Status(200).JSON(deeplinks)
}

func (h *Handler) GetDeeplink(c *fiber.Ctx) error {
	ctx := c.UserContext()

	slog.InfoContext(ctx, "Calling GetDeeplink in handler", slog.Any("test3", "testinfo3"))

	request := new(dto.GetDeeplinkRequest)
	if err := c.ParamsParser(request); err != nil {
		// TODO: using api error standard response
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	deeplink, err := h.deeplinkService.GetDeeplink(ctx, request)
	if err != nil {
		// TODO: using api error standard response
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// TODO: using api standard response
	return c.Status(200).JSON(deeplink)
}
