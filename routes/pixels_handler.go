package routes

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jonahgeorge/weatherglass/models"
	repo "github.com/jonahgeorge/weatherglass/repositories"
)

func (app *Application) PixelsCreateHandler(w http.ResponseWriter, r *http.Request) {
	eventsRepo := repo.NewEventsRepository(app.db)

	event := eventFromRequest(r)

	_, err := eventsRepo.Create(event)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func eventFromRequest(r *http.Request) *models.Event {
	queryParams, _ := ParseQuery(r.URL.RawQuery)

	siteId, _ := strconv.Atoi(queryParams.Get("site_id"))

	log.Println(*getOrNil(queryParams, "user_agent"))

	return &models.Event{
		SiteId:    siteId,
		Resource:  getOrNil(queryParams, "resource"),
		Referrer:  getOrNil(queryParams, "referrer"),
		Title:     getOrNil(queryParams, "title"),
		UserAgent: getOrNil(queryParams, "user_agent"),
	}
}

func getOrNil(values url.Values, key string) *string {
	val := values.Get(key)

	if len(val) > 0 {
		return &val
	} else {
		return nil
	}
}

func ParseQuery(query string) (url.Values, error) {
	m := make(url.Values)
	err := tolerantParseQuery(m, query)
	return m, err
}

func tolerantParseQuery(m url.Values, query string) (err error) {
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		m[key] = append(m[key], value)
	}
	return err
}
