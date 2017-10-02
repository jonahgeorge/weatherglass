package models

import "time"

type Event struct {
	Id             int
	SiteId         int
	Resource       *string
	Referrer       *string
	Title          *string
	UserAgent      *string
	CreatedAt      time.Time
	BrowserName    *string
	BrowserVersion *string
}

func (e Event) FromRow(row Scannable) error {
	return row.Scan(&e.Id, &e.SiteId, &e.Resource, &e.Referrer, &e.Title,
		&e.UserAgent, &e.CreatedAt, &e.BrowserName, &e.BrowserVersion)
}
