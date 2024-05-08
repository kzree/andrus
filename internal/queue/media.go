package queue

import "github.com/google/uuid"

type Media struct {
	id        string
	URL       string
	Title     string
	Requester *Requester
}

func NewMedia(url, title string, r *Requester) *Media {
	id := uuid.New().String()
	return &Media{URL: url, Title: title, Requester: r, id: id}
}
