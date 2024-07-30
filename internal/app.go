package internal

import (
	"cmp"
	"embed"
	"net/http"
	"os"

	"stripe-cop/internal/charge"
	"stripe-cop/internal/element"
	"stripe-cop/internal/home"
	"stripe-cop/internal/paymentintent"
	"stripe-cop/public"

	"github.com/leapkit/leapkit/core/assets"
	"github.com/leapkit/leapkit/core/db"
	"github.com/leapkit/leapkit/core/render"
	"github.com/leapkit/leapkit/core/server"
)

var (
	//go:embed **/*.html **/*.html *.html
	tmpls embed.FS

	// DB is the database connection builder function
	// that will be used by the application based on the driver and
	// connection string.
	DB = db.ConnectionFn(
		cmp.Or(os.Getenv("DATABASE_URL"), "leapkit.db"),
		db.WithDriver("sqlite3"),
	)
)

// Server interface exposes the methods
// needed to start the server in the cmd/app package
type Server interface {
	Addr() string
	Handler() http.Handler
}

func New() Server {
	// Creating a new server instance with the
	// default host and port values.
	r := server.New(
		server.WithHost(cmp.Or(os.Getenv("HOST"), "0.0.0.0")),
		server.WithPort(cmp.Or(os.Getenv("PORT"), "3000")),
		server.WithSession(
			cmp.Or(os.Getenv("SESSION_SECRET"), "d720c059-9664-4980-8169-1158e167ae57"),
			cmp.Or(os.Getenv("SESSION_NAME"), "leapkit_session"),
		),
	)

	assetsManager := assets.NewManager(public.Files)
	r.Use(render.Middleware(
		render.TemplateFS(tmpls, "internal"),

		render.WithDefaultLayout("layout.html"),
		render.WithHelpers(render.AllHelpers),
		render.WithHelpers(map[string]any{
			"assetPath": assetsManager.PathFor,
		}),
	))

	r.HandleFunc("GET /{$}", home.Index)

	r.HandleFunc("POST /create-charge/{$}", charge.PayChargeOne)
	r.HandleFunc("GET /back-ex-1/{$}", charge.BackOne)

	r.HandleFunc("POST /create-payment-confirmed/{$}", paymentintent.PayIntentConfirm)
	r.HandleFunc("GET /back-ex-2/{$}", paymentintent.BackTwo)

	r.HandleFunc("POST /create-payment-intent/{$}", element.HandleCreatePaymentIntent)

	r.HandleFunc("POST /create-charge-app-fee/{$}", charge.PayChargeWithAppFee)
	r.HandleFunc("GET /back-ex-4/{$}", charge.BackFour)

	r.HandleFunc("POST /create-payment-intent-app-fee/{$}", element.HandleCreatePaymentIntentAppFee)

	// Mounting the assets manager at the end of the routes
	// so that it can serve the public assets.
	r.HandleFunc(assetsManager.HandlerPattern(), assetsManager.HandlerFn)

	return r
}
