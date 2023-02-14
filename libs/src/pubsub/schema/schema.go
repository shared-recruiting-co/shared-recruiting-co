package schema

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	// Data is the message payload.
	Data []byte `json:"data"`
}

// EmailPushNotification is the payload of a gmail push notification event.
// Note the fields are camelCase, not snake_case, because Google generates them.
type EmailPushNotification struct {
	Email     string `json:"emailAddress"`
	HistoryID uint64 `json:"historyId"`
}

// EmailMessages is the payload for messages to be processes
type EmailMessages struct {
	// Email is the email address of the user.
	Email string `json:"email"`
	// Messages is the list of messages IDs.
	Messages []string `json:"messages"`
	// Settings allow us to configure the downstream processing behavior
	Settings EmailMessagesSettings `json:"settings"`
}

type EmailMessagesSettings struct {
	// DryRun indicates whether we should take action on the messages or not.
	DryRun bool `json:"dry_run"`
	// Reclassify indicates whether we should reclassify already classified messages or not.
	Reclassify bool `json:"reclassify"`
}
