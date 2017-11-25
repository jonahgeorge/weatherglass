package routes

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	"github.com/haisum/recaptcha"
	"github.com/jonahgeorge/weatherglass/models"
	repo "github.com/jonahgeorge/weatherglass/repositories"
	_ "github.com/lib/pq"
	"github.com/sendgrid/sendgrid-go"
)

type Application struct {
	db              *sql.DB
	sessions        *sessions.CookieStore
	emailClient     *sendgrid.Client
	recaptchaClient recaptcha.R
	hostName        string
}

func NewApplication() *Application {
	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	sessions := sessions.NewCookieStore([]byte(os.Getenv("SECRET_TOKEN")))
	emailClient := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	recaptchaClient := recaptcha.R{Secret: os.Getenv("RECAPTCHA_SECRET_TOKEN")}
	hostName := os.Getenv("HOST")

	return &Application{
		db:              db,
		sessions:        sessions,
		emailClient:     emailClient,
		recaptchaClient: recaptchaClient,
		hostName:        hostName,
	}
}

func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, data pongo2.Context) error {
	t, _ := pongo2.FromFile("templates/" + name + ".html")
	userRepo := repo.NewUsersRepository(app.db)
	session, _ := app.GetSession(r)

	if session.Values["userId"] != nil {
		userResult := <-userRepo.FindById(session.Values["userId"].(int))
		data["currentUser"] = userResult.Ok
	}

	data["flashes"] = session.Flashes()
	data["host"] = os.Getenv("HOST")
	data["recaptcha_site_key"] = os.Getenv("RECAPTCHA_SITE_KEY")
	data["weatherglass_site_id"] = os.Getenv("WEATHERGLASS_SITE_ID")
	session.Save(r, w)

	return t.ExecuteWriter(data, w)
}

func (app *Application) GetSession(r *http.Request) (*sessions.Session, error) {
	return app.sessions.Get(r, "weatherglass")
}

type AuthenticatedHandlerFunc func(http.ResponseWriter, *http.Request, *models.User)

func (app *Application) RequireAuthentication(next AuthenticatedHandlerFunc) http.HandlerFunc {
	userRepo := repo.NewUsersRepository(app.db)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := app.GetSession(r)
		if err != nil {
			session.AddFlash("You must be logged in!")
			session.Save(r, w)
			http.Redirect(w, r, "/login", 307)
			return
		}

		userId, ok := session.Values["userId"]
		if !ok {
			session.AddFlash("You must be logged in!")
			session.Save(r, w)
			http.Redirect(w, r, "/login", 307)
			return
		}

		userResult := <-userRepo.FindById(userId.(int))
		if userResult.Ok == nil || userResult.Err != nil {
			session.AddFlash("You must be logged in!")
			session.Save(r, w)
			http.Redirect(w, r, "/login", 307)
			return
		}

		next(w, r, userResult.Ok)
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
