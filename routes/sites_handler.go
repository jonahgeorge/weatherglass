package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/jonahgeorge/weatherglass/models"
	"github.com/jonahgeorge/weatherglass/queries"
	repo "github.com/jonahgeorge/weatherglass/repositories"
)

func (app *Application) SitesIndexHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	sites, _ := repo.NewSitesRepository(app.db).FindByUserId(currentUser.Id)

	app.Render(w, r, "sites/index", pongo2.Context{"sites": sites})
}

func (app *Application) SitesNewHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	app.Render(w, r, "sites/new", pongo2.Context{})
}

func (app *Application) SitesCreateHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	session, _ := app.GetSession(r)

	site := new(models.Site)
	site.Name = r.PostFormValue("name")

	_, err := repo.NewSitesRepository(app.db).Create(currentUser.Id, site.Name)
	if err != nil {
		session.AddFlash("An error occured while creating your site")
		session.Save(r, w)
		app.Render(w, r, "sites/new", pongo2.Context{"site": site})
		return
	}

	session.AddFlash("Successfully created site!")
	session.Save(r, w)

	http.Redirect(w, r, "/sites", 302)
}

func (app *Application) SitesShowHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	session, _ := app.GetSession(r)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	format := "2006-01-02 15:04:05"
	ending := time.Now()

	starting, granularity := parseRange(r.URL.Query().Get("range"))

	site, _ := repo.NewSitesRepository(app.db).FindById(id)
	if !currentUser.CanView(site) {
		session.AddFlash("You are not authorized to access this resource.")
		session.Save(r, w)
		http.Redirect(w, r, "/sites", 302)
		return
	}

	referrers, err := queries.NewReferrersQuery(app.db).Run(site.Id, starting, ending)
	if err != nil {
		log.Println(err)
	}

	userAgents, err := queries.NewUserAgentsQuery(app.db).Run(site.Id, starting, ending)
	if err != nil {
		log.Println(err)
	}

	app.Render(w, r, "sites/show", pongo2.Context{
		"site":        site,
		"starting":    starting.Format(format),
		"ending":      ending.Format(format),
		"granularity": granularity,
		"range":       r.URL.Query().Get("range"),
		"referrers":   referrers,
		"userAgents":  userAgents,
	})
}

func (app *Application) SitesEditHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	session, _ := app.GetSession(r)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	site, _ := repo.NewSitesRepository(app.db).FindById(id)
	if !currentUser.CanUpdate(site) {
		session.AddFlash("You are not authorized to access this resource.")
		session.Save(r, w)
		http.Redirect(w, r, "/sites", 302)
		return
	}

	app.Render(w, r, "sites/edit", pongo2.Context{"site": site})
}

func (app *Application) SitesUpdateHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	session, _ := app.GetSession(r)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	site, err := repo.NewSitesRepository(app.db).FindById(id)
	if !currentUser.CanUpdate(site) {
		session.AddFlash("You are not authorized to access this resource.")
		session.Save(r, w)
		http.Redirect(w, r, "/sites", 302)
		return
	}

	site.Name = r.PostFormValue("name")

	_, err = repo.NewSitesRepository(app.db).Update(id, site.Name)
	if err != nil {
		session.AddFlash("An error occured while updating this site")
		session.Save(r, w)
		app.Render(w, r, "sites/edit", pongo2.Context{"site": site})
		return
	}

	session.AddFlash("Successfully updated site!")
	session.Save(r, w)

	http.Redirect(w, r, fmt.Sprintf("/sites/%d", site.Id), 302)
}

func (app *Application) SitesDestroyHandler(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
	session, _ := app.GetSession(r)
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	site, _ := repo.NewSitesRepository(app.db).FindById(id)
	if !currentUser.CanDelete(site) {
		session.AddFlash("You are not authorized to access this resource.")
		session.Save(r, w)
		http.Redirect(w, r, "/sites", 302)
		return
	}

	_, err := repo.NewSitesRepository(app.db).Delete(site.Id)
	if err != nil {
		session.AddFlash("An error occured while deleting this site")
		session.Save(r, w)
		app.Render(w, r, "sites/edit", pongo2.Context{"site": site})
		return
	}

	session.AddFlash("Successfully deleted site!")
	session.Save(r, w)

	http.Redirect(w, r, "/sites", 302)
}

func parseRange(r string) (time.Time, string) {
	switch r {
	case "past-1-year":
		return time.Now().AddDate(-1, 0, 0), "30-day"
	case "past-1-month":
		return time.Now().AddDate(0, -1, 0), "1-day"
	case "past-7-days":
		return time.Now().Add(-7 * 24 * time.Hour), "2-hour"
	case "past-72-hours":
		return time.Now().Add(-74 * time.Hour), "1-hour"
	case "past-24-hours":
		return time.Now().AddDate(0, 0, -1), "10-minute"
	case "past-2-hours":
		return time.Now().Add(-2 * time.Hour), "1-minute"
	}

	return time.Now().Add(-2 * time.Hour), "1-minute"
}
