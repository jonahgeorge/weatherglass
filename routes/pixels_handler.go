package routes

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"strconv"

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
	queryParams := r.URL.Query()

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
