package honeycombio

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Queries describe all the query-related methods that the Honeycomb API
// supports.
//
// API docs: https://docs.honeycomb.io/api/queries/
type Queries interface {
	// Get a query by its ID.
	Get(ctx context.Context, dataset string, id string) (*QuerySpec, error)

	// Create a new query in this dataset. When creating a new query ID may
	// not be set.
	Create(ctx context.Context, dataset string, c *QuerySpec) (*QuerySpec, error)

	GetQueryResult(ctx context.Context, dataset string, queryID string) (*QueryResult, error)
}

// queries implements Queries.
type queries struct {
	client *Client
}

// Compile-time proof of interface implementation by type queries.
var _ Queries = (*queries)(nil)

func (s *queries) Get(ctx context.Context, dataset string, id string) (*QuerySpec, error) {
	var q QuerySpec
	err := s.client.performRequest(ctx, "GET", fmt.Sprintf("/1/queries/%s/%s", urlEncodeDataset(dataset), id), nil, &q)
	return &q, err
}

func (s *queries) Create(ctx context.Context, dataset string, data *QuerySpec) (*QuerySpec, error) {
	var q QuerySpec
	err := s.client.performRequest(ctx, "POST", "/1/queries/"+urlEncodeDataset(dataset), data, &q)
	return &q, err
}

type QueryResult struct {
	ID       string `json:"id"`
	Complete bool   `json:"complete"`
	Data     struct {
		Series []struct {
			Time time.Time              `json:"time"`
			Data map[string]interface{} `json:"data"`
		} `json:"series"`
	} `json:"data"`
}

func (s *queries) GetQueryResult(ctx context.Context, dataset string, queryID string) (*QueryResult, error) {
	q := struct {
		QueryID string `json:"query_id"`
	}{
		QueryID: queryID,
	}
	r := struct {
		ID       string `json:"id"`
		Complete bool   `json:"complete"`
	}{}

	err := s.client.performRequest(ctx, "POST", fmt.Sprintf("/1/query_results/%s", urlEncodeDataset(dataset)), q, &r)
	if err != nil {
		return nil, err
	}

	// Because query results are async, we may need to wait awhile for them to be available.
	qr := &QueryResult{}
	// Start by checking immediately, wait for 10 milliseconds the first time and double the wait time each time.
	var sleep = time.Millisecond * 10
	for i := 0; i < 10; i++ {
		err = s.client.performRequest(ctx, "GET", fmt.Sprintf("/1/query_results/%s/%s", urlEncodeDataset(dataset), r.ID), nil, qr)
		if err != nil {
			return nil, err
		}
		if qr.Complete {
			// We have the data, no need to keep checking.
			return qr, nil
		}

		// Data is not ready, sleep and try again.
		time.Sleep(sleep)
		sleep = sleep * 2
	}

	return nil, errors.New("Query timed out")
}
