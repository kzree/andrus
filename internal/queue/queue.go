package queue

import (
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Requester struct {
	ID       string
	Username string
}

type Media struct {
	id        string
	URL       string
	Title     string
	Requester *Requester
}

type Queue struct {
	logger   *zerolog.Logger
	maxItems int
	items    []*Media
}

func New(maxItems int, l *zerolog.Logger) *Queue {
	return &Queue{logger: l, maxItems: maxItems, items: make([]*Media, 0)}
}

func (q *Queue) Add(m *Media) error {
	q.logger.Info().Interface("media", m).Msg("adding media to queue")
	if len(q.items) < q.maxItems {
		q.logger.Debug().Int("size", len(q.items)).Int("max", q.maxItems).Msg("queue size before")
		id := uuid.New().String()
		q.logger.Debug().Str("id", id).Interface("media", m).Msg("generated id for media")
		m.id = id
		q.items = append(q.items, m)
		q.logger.Debug().Int("size", len(q.items)).Int("max", q.maxItems).Msg("queue size after")
		return nil
	}

	return errors.New("queue is full")
}

func (q *Queue) GetFirst() (*Media, error) {
	q.logger.Info().Msg("getting first media from queue")
	if len(q.items) == 0 {
		q.logger.Warn().Msg("queue is empty")
		return nil, errors.New("queue is empty")
	}

	q.logger.Debug().Int("size", len(q.items)).Int("max", q.maxItems).Msg("queue size before")
	first := q.items[0]
	q.items = q.items[1:]
	q.logger.Debug().Int("size", len(q.items)).Int("max", q.maxItems).Msg("queue size after")
	return first, nil
}

func (q *Queue) GetSize() int {
	return len(q.items)
}
