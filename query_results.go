package honeycombio

import (
	"context"
	"encoding/json"
	"fmt"
)

// QueryResults describe all the query results-related methods that the Honeycomb API
// supports.
//
// API docs: https://docs.honeycomb.io/api/query-results/
type QueryResults interface {
	// Get a query result by its ID.
	//
	// API docs: https://docs.honeycomb.io/api/query-results/#get-query-result
	Get(ctx context.Context, dataset string, id string) (*QueryResult, error)

	// Create a new result in this dataset. When creating a new query result ID may
	// not be set.
	//
	// API docs: https://docs.honeycomb.io/api/query-results/#create-query-result
	Create(ctx context.Context, dataset string, queryID string) (*QueryResult, error)
}

// queries implements Queries.
type queryResults struct {
	client *Client
}

// Compile-time proof of interface implementation by type queries.
var _ QueryResults = (*queryResults)(nil)

// QueryAnnotation represents a Honeycomb query result.
//
// API docs: https://docs.honeycomb.io/api/query-results/#get-example-response
type QueryResult struct {
	ID       string            `json:"id"`
	Complete bool              `json:"complete"`
	Data     *QueryResultData  `json:"data,omitempty"`
	Links    *QueryResultLinks `json:"links,omitempty"`
}

type QueryResultData struct {
	Series  json.RawMessage `json:"series"`
	Results json.RawMessage `json:"results"`
}

type QueryResultLinks struct {
	QueryURL      string `json:"query_url"`
	GraphImageURL string `json:"graph_image_url"`
}

func (s *queryResults) Get(ctx context.Context, dataset string, id string) (*QueryResult, error) {
	var q QueryResult
	err := s.client.performRequest(ctx, "GET", fmt.Sprintf("/1/query_results/%s/%s", urlEncodeDataset(dataset), id), nil, &q)
	return &q, err
}

func (s *queryResults) Create(ctx context.Context, dataset string, queryID string) (*QueryResult, error) {
	var q QueryResult

	data := map[string]string{
		"query_id": queryID,
	}
	err := s.client.performRequest(ctx, "POST", "/1/query_results/"+urlEncodeDataset(dataset), data, &q)
	return &q, err
}
