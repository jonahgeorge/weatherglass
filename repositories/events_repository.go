package repositories

import (
	"database/sql"

	"github.com/jonahgeorge/weatherglass/models"
)

const VISITS_ALL_SQL = "SELECT * FROM events"

const VISIT_CREATE_SQL = `
insert into events (site_id, resource, referrer, title, user_agent) 
values ($1, $2, $3, $4, $5)
returning *`

type EventsRepository struct {
	db *sql.DB
}

func NewEventsRepository(db *sql.DB) *EventsRepository {
	return &EventsRepository{db: db}
}

func (r EventsRepository) Create(e *models.Event) (*models.Event, error) {
	event := new(models.Event)
	row := r.db.QueryRow(VISIT_CREATE_SQL,
		e.SiteId, e.Resource, e.Referrer, e.Title, e.UserAgent)
	err := event.FromRow(row)
	return event, err
}
