package queue

import "github.com/google/uuid"

type Media struct {
	id        string
	URL       string
	Title     string
	Duration  string
	FilePath  *string
	Requester *Requester
}

func NewMedia(url, title string, duration string, r *Requester) *Media {
	id := uuid.New().String()
	return &Media{URL: url, Title: title, Duration: duration, Requester: r, id: id}
}

func (m *Media) SetMetadata(title string, duration string) {
	m.Title = title
	m.Duration = duration
}
