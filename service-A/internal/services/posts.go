package services

import (
	"context"
	"encoding/json"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type GetPosts struct {
	Success bool           `json:"success"`
	Payload GetPostPayload `json:"payload"`
}

type GetPostPayload struct {
	ID    string       `json:"id"`
	Title string       `json:"title"`
	Asset AssetPayload `json:"asset"`
}

type AssetPayload struct {
	ID      string `json:"id"`
	Key     string `json:"key"`
	AltText string `json:"alt_text"`
	URL     string `json:"url"`
}

func GetUsersPosts(ctx context.Context, postid string) (*GetPostPayload, error) {
	url := "http://localhost:8081/posts/" + postid

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body GetPosts
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &body.Payload, nil
}
