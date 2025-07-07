package ports

import (
	"context"
	"deeplink-bff/bff/internal/adapters/handler/dto"
)

type DeeplinkClient interface {
	GetDeeplinkList(ctx context.Context) (*dto.GetDeeplinkListResponse, error)
	GetDeeplink(ctx context.Context, id string) (*dto.GetDeeplinkResponse, error)
}
