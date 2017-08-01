package models

import "time"

type Event struct {
	Id        int
	SiteId    int
	Resource  *string
	Referrer  *string
	Title     *string
	UserAgent *string
	CreatedAt time.Time
}
