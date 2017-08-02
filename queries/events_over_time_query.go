package queries

import "database/sql"

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
    date_trunc($3, timestamp) as interval,
    count(*) as count
  from events
  where site_id = $1
  group by interval
) as events_per_interval on events_per_interval.interval = timeslices.interval
order by timeslices.interval`

type EventsOverTimeQuery struct {
	db *sql.DB
}

func NewEventsPerMinuteQuery(db *sql.DB) *EventsOverTimeQuery {
	return &EventsOverTimeQuery{db: db}
}

func (q *EventsOverTimeQuery) Run(siteId int, wange string, interval string) error {
	rows, err := q.db.Query(EVENTS_PER_MINUTE_SQL, siteId, wange, interval, "1 "+interval)

	return err
}
