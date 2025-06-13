package posthog

import (
	"onepixel_backend/src/config"
	"onepixel_backend/src/utils/applogger"
	"sync"

	"github.com/posthog/posthog-go"
)

var (
	client posthog.Client
	once   sync.Once
)

// GetPostHogClient returns a singleton PostHog client
func GetPostHogClient() posthog.Client {
	once.Do(func() {
		if config.PosthogApiKey == "" || config.PosthogApiKey == "xxxx" {
			applogger.Warn("PostHog API key not configured, PostHog events will be disabled")
			return
		}

		var err error
		client, err = posthog.NewWithConfig(
			config.PosthogApiKey,
			posthog.Config{
				Endpoint: "https://us.i.posthog.com",
			},
		)
		if err != nil {
			applogger.Error("Failed to initialize PostHog client: ", err)
			return
		}
		applogger.Info("PostHog client initialized successfully")
	})
	return client
}

// TrackEvent sends an event to PostHog
// This function is non-blocking and handles errors gracefully
func TrackEvent(distinctId string, event string, properties map[string]interface{}) {
	client := GetPostHogClient()
	if client == nil {
		// Client not initialized (likely due to missing API key)
		return
	}

	err := client.Enqueue(posthog.Capture{
		DistinctId: distinctId,
		Event:      event,
		Properties: properties,
	})

	if err != nil {
		applogger.Error("Failed to send PostHog event: ", err)
	}
}

// Close closes the PostHog client
func Close() {
	if client != nil {
		if err := client.Close(); err != nil {
			applogger.Error("Failed to close PostHog client: ", err)
		}
	}
}
