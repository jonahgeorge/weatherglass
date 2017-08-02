package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jonahgeorge/weatherglass/models"
	"github.com/jonahgeorge/weatherglass/queries"
)

func (app *Application) EventsOverTimeIndexHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	siteId, _ := strconv.Atoi(mux.Vars(r)["id"])
	interval := r.URL.Query().Get("interval")
	wange := r.URL.Query().Get("range")

	eventsPerMinuteQuery := queries.NewEventsPerMinuteQuery(app.db)
	eventsPerMinuteQuery.Run(siteId, wange, interval)

	// const results = {};
	// report.forEach((r) => { results[r["interval"]] = r["count"]; });
	// response.json(results);
}
