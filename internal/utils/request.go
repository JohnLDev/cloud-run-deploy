package utils

import (
	"context"
	"errors"
	"io"
	"net/http"
	"slices"
)

func RequestWithContext(ctx context.Context, url string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	result, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	response, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	if !slices.Contains([]int{http.StatusOK, http.StatusAccepted, http.StatusCreated}, result.StatusCode) {
		return nil, errors.New(string(response))
	}

	return response, nil
}
