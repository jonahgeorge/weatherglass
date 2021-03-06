package queries

import (
	"database/sql"
	"log"
	"time"
)

const REFERRERS_SQL = `
select distinct 
  host, 
  count(*)
from (
  select
    substring(referrer from '.*://([^/]*)') as host
  from events
  where site_id = $1
    and created_at between $2 and $3
    and referrer is not null
) referrers
group by host
order by count desc`

type ReferrersQuery struct {
	db *sql.DB
}

type ReferrersResult struct {
	Referrer string `json:"referrer"`
	Count    int    `json:"count"`
}

func NewReferrersQuery(db *sql.DB) *ReferrersQuery {
	return &ReferrersQuery{db: db}
}

func (q *ReferrersQuery) Run(siteId int, starting time.Time, ending time.Time) ([]ReferrersResult, error) {
	var results []ReferrersResult
	rows, err := q.db.Query(REFERRERS_SQL, siteId, starting, ending)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		result := new(ReferrersResult)
		err = rows.Scan(&result.Referrer, &result.Count)
		results = append(results, *result)
	}

	return results, err
}
