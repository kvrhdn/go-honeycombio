package honeycombio

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type eventTestInput struct {
	Column1    string  `json:"column_1"`
	DurationMs float64 `json:"duration_ms"`
}

func TestEvents(t *testing.T) {
	ctx := context.Background()

	c := newTestClient(t)
	dataset := testDataset(t)

	t.Run("Send", func(t *testing.T) {
		data := eventTestInput{
			Column1:    "foo",
			DurationMs: 1000,
		}

		err := c.Events.Send(ctx, dataset, data)

		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("SendBatch", func(t *testing.T) {
		data := []SendBatchRequest{
			{
				Data: eventTestInput{
					Column1:    "foo",
					DurationMs: 1000,
				},
			},
			{
				Time:       TimePtr(time.Now()),
				SampleRate: 2,
				Data: eventTestInput{
					Column1:    "bar",
					DurationMs: 2000,
				},
			},
		}

		responses, err := c.Events.SendBatch(ctx, dataset, data)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, len(data), len(responses))
		for _, response := range responses {
			assert.Equal(t, 202, response.Status)
		}
	})
}
