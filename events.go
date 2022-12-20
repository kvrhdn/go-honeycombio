package honeycombio

import (
	"context"
	"time"
)

// Events describe all the events-related methods that the Honeycomb API
// supports.
//
// API docs: https://docs.honeycomb.io/api/events/
type Events interface {
	// Send a single event to this dataset.
	Send(ctx context.Context, dataset string, data map[string]interface{}) error

	// SendBatch sends batch of events to this dataset
	SendBatch(ctx context.Context, dataset string, data []SendBatchRequest) ([]SendBatchResponse, error)
}

// events implements Events.
type events struct {
	client *Client
}

// Compile-time proof of interface implementation by type queries.
var _ Events = (*events)(nil)

// SendBatchRequest represents event batch request body
type SendBatchRequest struct {
	Time       *time.Time             `json:"time,omitempty"`
	SampleRate int                    `json:"samplerate,omitempty"`
	Data       map[string]interface{} `json:"data"`
}

// SendBatchRequest represents event batch response body
type SendBatchResponse struct {
	Status int `json:"status"`
}

// Send a single event to this dataset.
func (s *events) Send(ctx context.Context, dataset string, data map[string]interface{}) error {
	return s.client.performRequest(ctx, "POST", "/1/events/"+urlEncodeDataset(dataset), data, nil)
}

// SendBatch sends batch of events to this dataset
func (s *events) SendBatch(ctx context.Context, dataset string, data []SendBatchRequest) ([]SendBatchResponse, error) {
	var q []SendBatchResponse
	err := s.client.performRequest(ctx, "POST", "/1/batch/"+urlEncodeDataset(dataset), data, &q)
	return q, err
}
