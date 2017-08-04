package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jonahgeorge/weatherglass/models"
	repo "github.com/jonahgeorge/weatherglass/repositories"
	_ "github.com/lib/pq"
)

func main() {
	app := NewApplication()

	r := mux.NewRouter()
	r.HandleFunc("/", app.IndexHandler).Methods("GET")
	r.HandleFunc("/login", app.SessionsNewHandler).Methods("GET")
	r.HandleFunc("/login", app.SessionsCreateHandler).Methods("POST")
	r.HandleFunc("/logout", app.SessionsDestroyHandler).Methods("GET")
	r.HandleFunc("/signup", app.UsersNewHandler).Methods("GET")
	r.HandleFunc("/signup", app.UsersCreateHandler).Methods("POST")
	r.HandleFunc("/email_confirmation/new", app.EmailConfirmationsNewHandler).Methods("GET")
	r.HandleFunc("/email_confirmation", app.EmailConfirmationsCreateHandler).Methods("POST")
	r.HandleFunc("/email_confirmation", app.EmailConfirmationsShowHandler).Methods("GET")
	r.HandleFunc("/track.gif", app.PixelsCreateHandler).Methods("GET")
	r.HandleFunc("/documentation", app.DocumentationIndexHandler).Methods("GET")

	r.HandleFunc("/sites", app.RequireAuthentication(app.RequireEmailConfirmation(app.SitesIndexHandler))).Methods("GET")
	r.HandleFunc("/sites/{id:[0-9]+}", app.RequireAuthentication(app.RequireEmailConfirmation(app.SitesShowHandler))).Methods("GET")
	r.HandleFunc("/sites/{id:[0-9]+}/edit", app.RequireAuthentication(app.RequireEmailConfirmation(app.SitesEditHandler))).Methods("GET")
	r.HandleFunc("/sites/{id:[0-9]+}", app.RequireAuthentication(app.RequireEmailConfirmation(app.SitesUpdateHandler))).Methods("PUT")
	r.HandleFunc("/sites/{id:[0-9]+}", app.RequireAuthentication(app.RequireEmailConfirmation(app.SitesDestroyHandler))).Methods("DELETE")
	r.HandleFunc("/sites/new", app.RequireAuthentication(app.RequireEmailConfirmation(app.SitesNewHandler))).Methods("GET")
	r.HandleFunc("/sites", app.RequireAuthentication(app.RequireEmailConfirmation(app.SitesCreateHandler))).Methods("POST")
	r.HandleFunc("/sites/{id:[0-9]+}/reports/events_over_time", app.RequireAuthentication(app.RequireEmailConfirmation(app.EventsOverTimeIndexHandler))).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}

	log.Println("Listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port,
		handlers.CompressHandler(
			handlers.HTTPMethodOverrideHandler(
				handlers.LoggingHandler(os.Stdout, r)))))
}

type AuthenticatedHandlerFunc func(http.ResponseWriter, *http.Request, *models.User)

func (app *Application) RequireAuthentication(next AuthenticatedHandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.GetSession(r)

		userId, ok := session.Values["userId"]
		if !ok {
			session.AddFlash("You must be logged in!")
			session.Save(r, w)
			http.Redirect(w, r, "/login", 307)
			return
		}

		user, err := repo.NewUsersRepository(app.db).FindById(userId.(int))
		if user == nil || err != nil {
			session.AddFlash("You must be logged in!")
			session.Save(r, w)
			http.Redirect(w, r, "/login", 307)
			return
		}

		next(w, r, user)
	})
}

func (app *Application) RequireEmailConfirmation(next AuthenticatedHandlerFunc) AuthenticatedHandlerFunc {
	return AuthenticatedHandlerFunc(func(w http.ResponseWriter, r *http.Request, currentUser *models.User) {
		session, _ := app.GetSession(r)

		if !currentUser.IsEmailConfirmed {
			session.AddFlash("You must confirm your email address before continuing")
			session.Save(r, w)
			http.Redirect(w, r, "/email_confirmation/new", 302)
			return
		}

		next(w, r, currentUser)
	})
}
