package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// NotificationsClient sends events to the corporate notifications-api.
// Implements handler.NotificationPublisher interface.
type NotificationsClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// NewNotificationsClient creates a new notifications client.
// If baseURL is empty, PublishEvent becomes a no-op (dev mode).
func NewNotificationsClient(baseURL, token string) *NotificationsClient {
	return &NotificationsClient{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// PublishEvent sends an event to the notifications-api.
// No-op if baseURL is empty.
func (c *NotificationsClient) PublishEvent(event map[string]interface{}) error {
	if c.baseURL == "" {
		return nil // dev mode, no notifications
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("notifications: marshal event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/internal/v1/events",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("notifications: create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("notifications: send event: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("notifications: API returned %d", resp.StatusCode)
	}

	return nil
}
