package queries

import (
	"database/sql"
	"log"
	"time"
)

const userAgentsSql = `
select distinct 
  user_agent, 
  count(*)
from events
where site_id = $1
  and created_at between $2 and $3
  and user_agent is not null
group by user_agent
order by count desc`

type UserAgentsQuery struct {
	db *sql.DB
}

type UserAgentsResult struct {
	UserAgent string `json:"user_agent"`
	Count     int    `json:"count"`
}

func NewUserAgentsQuery(db *sql.DB) *UserAgentsQuery {
	return &UserAgentsQuery{db: db}
}

func (q *UserAgentsQuery) Run(siteId int, starting time.Time, ending time.Time) ([]UserAgentsResult, error) {
	var results []UserAgentsResult
	rows, err := q.db.Query(userAgentsSql, siteId, starting, ending)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		result := new(UserAgentsResult)
		err = rows.Scan(&result.UserAgent, &result.Count)
		results = append(results, *result)
	}

	return results, err
}
