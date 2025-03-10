package core

import (
	// "context"

	"github.com/oklog/ulid/v2"
	// "github.com/sashabaranov/go-openai"
)

type History interface {
	SaveMessage() ulid.ULID
}

// Target represents an LLM completion endpoint and model to run on it.
type Target struct {
	Endpoint string
	Model    string
	Key      string
}

// Request represents a completion request sent to a Target.
type Request struct {
	Id        ulid.ULID
	MessageId ulid.ULID
}

// Response represents a completion response returned by a Target for a Request.
type Response struct {
	Id        ulid.ULID
	MessageId ulid.ULID
}

// ParseResponseMessage parses m and saves it into h, returning a unique
// Response with an Id that is (should be) known to h.
func ParseResponseMessage(m string) {}

// Roundtrip represents an element in a Conversation.
type Roundtrip struct {
	Id       ulid.ULID
	Target   *Target
	Request  *Request
	Response *Response
}

// Conversation represents a series of Roundtrips.
type Conversation struct {
}
