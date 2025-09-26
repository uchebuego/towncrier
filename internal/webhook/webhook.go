package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"net/http"
	"time"
)

func sanitizeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "<invalid-url>"
	}
	return fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.Path)
}

func Send(webhookURL string, data map[string]interface{}, eventName string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Webhook [%s]: failed to marshal event payload: %v", eventName, err)
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	safeURL := sanitizeURL(webhookURL)
	log.Printf("Webhook [%s]: dispatching to %s with payload: %s", eventName, safeURL, string(jsonData))

	maxAttempts := 5
	baseDelay := 500 * time.Millisecond
	client := &http.Client{Timeout: 10 * time.Second}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		resp, err := client.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
		if err == nil {
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				resp.Body.Close()
				log.Printf("Webhook [%s]: success (status %d) on attempt %d", eventName, resp.StatusCode, attempt)
				return nil
			}
			statusCode := resp.StatusCode
			resp.Body.Close()
			err = fmt.Errorf("received non-2xx response: %v", statusCode)
		}

		log.Printf("Webhook [%s]: attempt %d failed: %v", eventName, attempt, err)

		if attempt == maxAttempts {
			log.Printf("Webhook [%s]: giving up after %d attempts: %v", eventName, attempt, err)
			return fmt.Errorf("webhook POST failed after %d attempts: %w", attempt, err)
		}

		sleep := baseDelay * time.Duration(1<<uint(attempt-1))
		time.Sleep(sleep)
	}

	return nil
}
