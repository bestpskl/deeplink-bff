package deeplink_service

import (
	"context"
	"deeplink-bff/bff/internal/adapters/handler/dto"
	"deeplink-bff/bff/internal/core/ports"
	"errors"
	"log/slog"
)

type deeplinkService struct {
	deeplinkClient ports.DeeplinkClient
}

func NewDeeplinkService(deeplinkClient ports.DeeplinkClient) ports.DeeplinkService {
	return &deeplinkService{
		deeplinkClient,
	}
}

func (d *deeplinkService) GetDeeplinkList(ctx context.Context) (*dto.GetDeeplinkListResponse, error) {

	jsonData := `{
			"user": {
					"email": "john@example.com",
					"password": "secret123",
					"preferences": {
							"theme": "testinfo2",
							"api_key": "key123"
					}
			}
	}`
	slog.InfoContext(ctx, "Calling GetDeeplinkList in service", slog.Any("user", jsonData))

	// deeplinks, err := d.deeplinkClient.GetDeeplinkList(ctx)

	// return deeplinks, err

	deeplinks := &dto.GetDeeplinkListResponse{
		Deeplinks: []dto.GetDeeplinkResponse{
			{
				Email: "john@example.com",
			},
		},
	}

	return deeplinks, nil
}

func (d *deeplinkService) GetDeeplink(ctx context.Context, request *dto.GetDeeplinkRequest) (*dto.GetDeeplinkResponse, error) {

	slog.InfoContext(ctx, "Calling GetDeeplink in service", slog.Any("deeplink", "testinfo4"))
	// deeplink, err := d.deeplinkClient.GetDeeplink(ctx, request.Id)
	slog.ErrorContext(ctx, "Calling GetDeeplink in service failed", slog.Any("error", "testerror4"))

	// return deeplink, err

	err := errors.New("unexpect error")

	return nil, err
}
