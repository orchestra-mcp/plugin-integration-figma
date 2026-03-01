package figma

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const baseURL = "https://api.figma.com/v1"

type Client struct {
	Token string
}

func NewClient() *Client {
	return &Client{Token: os.Getenv("FIGMA_ACCESS_TOKEN")}
}

func (c *Client) Get(ctx context.Context, path string) ([]byte, error) {
	if c.Token == "" {
		return nil, fmt.Errorf("FIGMA_ACCESS_TOKEN not set")
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", baseURL+"/"+path, nil)
	req.Header.Set("X-Figma-Token", c.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Figma API error %d: %s", resp.StatusCode, body)
	}
	return body, nil
}

func (c *Client) GetFormatted(ctx context.Context, path string) (string, error) {
	body, err := c.Get(ctx, path)
	if err != nil {
		return "", err
	}
	var v any
	json.Unmarshal(body, &v)
	formatted, _ := json.MarshalIndent(v, "", "  ")
	return string(formatted), nil
}
