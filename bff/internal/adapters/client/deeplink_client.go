package client

import (
	"context"
	"deeplink-bff/bff/internal/adapters/handler/dto"
	"encoding/json"
	"fmt"
	"net/http"
)

type DeeplinkClient struct {
	baseUrl string
}

func NewDeepLinkClient(baseUrl string) *DeeplinkClient {
	return &DeeplinkClient{
		baseUrl,
	}
}

func (d *DeeplinkClient) GetDeeplinkList(ctx context.Context) (*dto.GetDeeplinkListResponse, error) {

	url := fmt.Sprintf("%s/api/v1/deeplink", d.baseUrl)

	fmt.Println("--------------------")
	fmt.Println(url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	webclientResponse := new(dto.GetDeeplinkListResponse)
	err = json.NewDecoder(resp.Body).Decode(&webclientResponse)

	return webclientResponse, err
}

func (d *DeeplinkClient) GetDeeplink(ctx context.Context, id string) (*dto.GetDeeplinkResponse, error) {
	url := fmt.Sprintf("%s/api/v1/deeplink/%s", d.baseUrl, id)

	fmt.Println("--------------------")
	fmt.Println(url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	webclientResponse := new(dto.GetDeeplinkResponse)
	err = json.NewDecoder(resp.Body).Decode(&webclientResponse)

	return webclientResponse, err
}
