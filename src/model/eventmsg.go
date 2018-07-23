package model

// EventMessage is just a sample type
type EventMessage struct {
	// ID identifies this message.
	ID string `json:"id,omitempty"`

	// Labels represents the key-value pairs the entry is lebeled with.
	Labels map[string]string `json:"labels,omitempty"`

	// Message
	Data string `json:"data"`
}
