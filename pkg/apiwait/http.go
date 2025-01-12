package apiwait

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// HTTP waits for an HTTP server to become available at the given address.
// It polls /health until it gets a 200 OK response.
func HTTP(ctx context.Context, address string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	url := fmt.Sprintf("http://%s/healthcheck", address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	zap.L().Debug("checking health", zap.String("url", url))

	client := http.Client{}
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		zap.L().Debug(
			"health check",
			zap.String("url", url),
			zap.Int("status", resp.StatusCode),
		)
		_ = resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		if resp.StatusCode >= 500 {
			time.Sleep(time.Second)
			continue
		}
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
