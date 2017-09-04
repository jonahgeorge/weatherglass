package queries

import (
	"database/sql"
	"log"
	"strings"
	"time"
)

const EVENTS_PER_MINUTE_SQL = `
with timeslices as (
  select
    interval,
    0 as blank_count
  from generate_series(
    to_timestamp(floor(extract(epoch from $2::timestamp) / extract(epoch from $4::interval)) * extract(epoch from $4::interval)), 
    to_timestamp(floor(extract(epoch from $3::timestamp) / extract(epoch from $4::interval)) * extract(epoch from $4::interval)), 
    $4::interval 
  ) as interval
),
events_per_interval as (
  select
    to_timestamp(floor(extract(epoch from created_at) / extract(epoch from $4::interval)) * extract(epoch from $4::interval)) as interval, 
    count(*) as count
  from events
  where site_id = $1
    and created_at >= $2
  group by interval
)
select 
  timeslices.interval as interval, 
  coalesce(events_per_interval.count, timeslices.blank_count) as count
from timeslices
  left outer join events_per_interval on events_per_interval.interval = timeslices.interval
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

func splitLink(s, sep string) (string, string) {
	x := strings.Split(s, sep)
	return x[0], x[1]
}

func (q *EventsOverTimeQuery) Run(siteId int, starting string, ending string, granularity string) ([]EventsOverTimeResult, error) {
	granularity = strings.Replace(granularity, "-", " ", 1)

	var results []EventsOverTimeResult
	rows, err := q.db.Query(EVENTS_PER_MINUTE_SQL, siteId, starting, ending, granularity)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		result := new(EventsOverTimeResult)
		err = rows.Scan(&result.Interval, &result.Count)
		results = append(results, *result)
	}

	return results, err
}
