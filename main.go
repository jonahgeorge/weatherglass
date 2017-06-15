package main

import (
	"net/http"
	"log"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/julienschmidt/httprouter"
	"github.com/flosch/pongo2"
	"github.com/urfave/negroni"
	"github.com/goincremental/negroni-sessions"
	"github.com/goincremental/negroni-sessions/cookiestore"
	// "github.com/rs/cors"
	"os"
)

type Application struct {
	db *sql.DB
}

func NewApplication() *Application {
	dataSourceName := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	app := &Application{
		db: db,
	}

	return app
}

func (app *Application) RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data pongo2.Context) error {
	var t *pongo2.Template
	var err error

	t, err = pongo2.FromFile("templates/" + name + ".html")

	session := sessions.GetSession(r)

	if id, ok := session.Get("userId").(int); ok {
		user := NewUserRepository(app.db).FindById(id)
		data["currentUser"] = user
	}

	// Add some static values
	data["flashes"] = session.Flashes()

	if err != nil {
		return err
	}

	return t.ExecuteWriter(data, w)
}

func main() {
	app := NewApplication()

	store := cookiestore.New([]byte("keyboardcat"))

	router := httprouter.New()
	router.GET("/login", app.SessionsNewHandler)
	router.POST("/login", app.SessionsCreateHandler)
	router.GET("/logout", app.SessionsDestroyHandler)
	router.GET("/sites", app.SitesIndexHandler)
	router.GET("/documentation", app.DocumentationIndexHandler)
	router.GET("/", app.RootIndexHandler)

	n := negroni.Classic()
	n.Use(sessions.Sessions("weatherglass_session", store))
	n.UseHandler(router)
	n.Run()
}
