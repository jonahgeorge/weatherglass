package queries

import (
	"database/sql"
	"time"
)

const EVENTS_PER_MINUTE_SQL = `
with timeslices as (
  select
    interval,
    0 as blank_count
  from generate_series(
    date_trunc($2, (current_timestamp at time zone 'utc'))::timestamp,
    (current_timestamp at time zone 'utc')::timestamp, 
    $4
  ) as interval
)
select
  timeslices.interval as interval,
  coalesce(events_per_interval.count, timeslices.blank_count) as count
from timeslices
left outer join (
  select
    date_trunc($3, created_at) as interval,
    count(*) as count
  from events
  where site_id = $1
  group by interval
) as events_per_interval on events_per_interval.interval = timeslices.interval
order by timeslices.interval`

type EventsOverTimeQuery struct {
	db *sql.DB
}

type EventsOverTimeResult struct {
	Interval time.Time `json:"interval"`
	Count    int       `json:"count"`
}

func NewEventsPerMinuteQuery(db *sql.DB) *EventsOverTimeQuery {
	return &EventsOverTimeQuery{db: db}
}

func (q *EventsOverTimeQuery) Run(siteId int, wange string, interval string) ([]EventsOverTimeResult, error) {
	var results []EventsOverTimeResult
	rows, err := q.db.Query(EVENTS_PER_MINUTE_SQL, siteId, wange, interval, "1 "+interval)

	for rows.Next() {
		result := new(EventsOverTimeResult)
		err = rows.Scan(&result.Interval, &result.Count)
		results = append(results, *result)
	}

	return results, err
}
