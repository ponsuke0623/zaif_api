package zaif

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// BaseURL is an endpoint of api
const BaseURL = "https://api.zaif.jp/api/1/%s"

// Client interfaces is an interface for client
type Client interface {
	GetPairs(ctx context.Context) (Res, error)
	GetPrice(ctx context.Context, pair string) (Res, error)
	GetTicker(ctx context.Context, pair string) (Res, error)
	GetTrades(ctx context.Context, pair string) (Res, error)
}

type client struct{}

// New creates a new client
func New() Client {
	return &client{}
}

func (c *client) do(ctx context.Context, method, path string, body io.Reader) (json.RawMessage, error) {
	url := fmt.Sprintf(BaseURL, path)
	client := &http.Client{}

	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequest failed")
	}

	// set context
	r = r.WithContext(ctx)

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// status code is checked
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("some error is occurred with status code: %d", res.StatusCode)
	}

	var rawMsg json.RawMessage
	if err := json.NewDecoder(res.Body).Decode(&rawMsg); err != nil {
		return nil, err
	}
	return rawMsg, nil
}
