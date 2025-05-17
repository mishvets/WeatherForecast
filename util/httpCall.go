package util

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetRequest(ctx context.Context, url string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("get request fail: %w", err)
	}
	defer response.Body.Close()

	responseByte, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("api response: - %s, err - %s", responseByte, err) //TODO: delete

	return responseByte, nil
}
