package routes

import (
	"net/http"

	"github.com/flosch/pongo2"
)

func (app *Application) DocumentationIndexHandler(w http.ResponseWriter, r *http.Request) {
	app.Render(w, r, "documentation/index", pongo2.Context{})
}
