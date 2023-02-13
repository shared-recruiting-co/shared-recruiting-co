package cloudfunctions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/getsentry/sentry-go"

	"github.com/shared-recruiting-co/shared-recruiting-co/libs/src/pubsub/schema"
)

// const (
// provider = "google"
// )

func init() {
	functions.CloudEvent("Handle", Handle)
}

// naive error handler for now
func handleError(msg string, err error) error {
	err = fmt.Errorf("%s: %w", msg, err)
	sentry.CaptureException(err)
	return err
}

func Handle(ctx context.Context, e event.Event) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
		ServerName:       "email-push-notifications",
	})
	if err != nil {
		return fmt.Errorf("sentry.Init: %v", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.RecoverWithContext(ctx)

	var msg schema.MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return handleError("error parsing Pub/Sub message", err)
	}

	data := msg.Message.Data
	log.Printf("Event: %s", data)

	var payload schema.EmailMessages
	if err := json.Unmarshal(data, &payload); err != nil {
		return handleError("error parsing email messages", err)
	}

	log.Printf("Payload: %v", payload)

	return nil
}
