package queries

import (
	"database/sql"
	"log"
	"time"
)

const devicesSql = `
select distinct 
  browser_name,
  browser_version,
  count(*)
from events
where site_id = $1
  and created_at between $2 and $3
  and browser_name is not null
  and browser_version is not null
group by browser_name, browser_version
order by count desc`

type DevicesQuery struct {
	db *sql.DB
}

type DevicesResult struct {
	BrowserName    string
	BrowserVersion string
	Count          int
}

func NewDevicesQuery(db *sql.DB) *DevicesQuery {
	return &DevicesQuery{db: db}
}

func (q *DevicesQuery) Run(siteId int, starting time.Time, ending time.Time) ([]DevicesResult, error) {
	var results []DevicesResult
	rows, err := q.db.Query(devicesSql, siteId, starting, ending)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		result := new(DevicesResult)
		err = rows.Scan(&result.BrowserName, &result.BrowserVersion, &result.Count)
		results = append(results, *result)
	}

	return results, err
}
