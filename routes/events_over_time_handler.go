package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jonahgeorge/weatherglass/models"
	"github.com/jonahgeorge/weatherglass/queries"
)

func (app *Application) EventsOverTimeIndexHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	siteId, _ := strconv.Atoi(mux.Vars(r)["id"])

	eventsPerMinuteQuery := queries.NewEventsPerMinuteQuery(app.db)

	results, _ := eventsPerMinuteQuery.Run(
		siteId,
		r.URL.Query().Get("starting"),
		r.URL.Query().Get("ending"),
		r.URL.Query().Get("granularity"),
	)

	response := make(map[time.Time]int)
	for _, v := range results {
		response[v.Interval] = v.Count
	}

	b, _ := json.MarshalIndent(response, "", " ")

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
