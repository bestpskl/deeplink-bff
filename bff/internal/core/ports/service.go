package ports

import (
	"context"
	"deeplink-bff/bff/internal/adapters/handler/dto"
)

type DeeplinkService interface {
	GetDeeplinkList(ctx context.Context) (*dto.GetDeeplinkListResponse, error)
	GetDeeplink(ctx context.Context, request *dto.GetDeeplinkRequest) (*dto.GetDeeplinkResponse, error)
}
