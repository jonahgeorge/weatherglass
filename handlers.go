package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/flosch/pongo2"
	"github.com/goincremental/negroni-sessions"
)

func (app *Application) RootIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.RenderTemplate(w, r, "root/index", pongo2.Context{})
}

func (app *Application) DocumentationIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.RenderTemplate(w, r, "documentation/index", pongo2.Context{})
}

func (app *Application) SessionsNewHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	app.RenderTemplate(w, r, "sessions/new", pongo2.Context{})
}

func (app *Application) SessionsCreateHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessions.GetSession(r)
	userRepo := NewUserRepository(app.db)

	user := userRepo.FindByEmailAndPassword(r.PostFormValue("email"), r.PostFormValue("password"))

	if user != nil {
		session.Set("userId", user.id)
		session.AddFlash("Successfully logged in!")
		http.Redirect(w, r, "/", 302)
	} else {
		session.AddFlash("Either your email or password was invalid.")
		app.RenderTemplate(w, r, "sessions/new", pongo2.Context{
			"email": r.PostFormValue("email"),
		})
	}
}

func (app *Application) SessionsDestroyHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session := sessions.GetSession(r)
	session.Delete("userId")
	session.AddFlash("Successfully logged out!")

	http.Redirect(w, r, "/", 302)
}

func (app *Application) SitesIndexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sitesRepo := NewSitesRepository(app.db)

	app.RenderTemplate(w, r, "sites/index", pongo2.Context{
		"sites": sitesRepo.FindAll(),
	})
}
