package model

// EventClient is just a sample type
type EventClient struct {
	// ID identifies this message.
	ID string `json:"id,omitempty"`

	// Labels represents the key-value pairs the entry is lebeled with.
	Labels map[string]string `json:"labels,omitempty"`

	// Message
	Message chan EventMessage
}
