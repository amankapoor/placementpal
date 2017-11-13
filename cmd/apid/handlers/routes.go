// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package handlers

import (
	"log"
	"net/http"
	"sort"

	"github.com/amankapoor/goth"
	"github.com/amankapoor/goth/providers/gplus"
	"github.com/amankapoor/placementpal/internal/middleware"
	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/amankapoor/placementpal/internal/platform/web"
)

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

var ProviderIndexVar = &ProviderIndex{}

// API returns a handler for a set of routes.
func API(masterDB *db.DB) http.Handler {

	goth.UseProviders(
		gplus.New("1030362162437-qjd6mi9ld2gd0v2kiml6orvd590vpfjm.apps.googleusercontent.com", "rhEICw3dqRBgGDDcRXDirESm", "http://localhost:8080/auth/gplus/callback"),
	)

	var keys []string
	m := make(map[string]string)
	m["gplus"] = "Google Plus"

	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ProviderIndexVar = &ProviderIndex{Providers: keys, ProvidersMap: m}

	// Create the web handler for setting routes and middleware.
	app := web.New(middleware.ReqResLogger, middleware.ErrorHandler)
	//This is gated handler that will require login
	requiresLogin := app.Group(requireLogin)
	requiresAdmin := app.Group(requireAdmin)

	// cacheInvalid := app.Group(cacheInvalidator)

	// NotFound handlers displays custom page and logs the request
	app.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 %s -> %s : %s",
			r.URL.Path,
			r.RemoteAddr,
			r.UserAgent(),
		)
		w.WriteHeader(404)
		w.Write([]byte("This page is not available."))
	})

	// Serve static views directory
	app.Router.ServeFiles("/views/*filepath", http.Dir("views"))

	app.Router.ServeFiles("/files/*filepath", http.Dir("files"))

	db := Database{masterDB}

	app.Handle("GET", "/", Home)
	app.Handle("GET", "/auth/:provider", db.LoginWithGoogle)
	app.Handle("GET", "/auth/:provider/callback", db.LoginWithGoogleCallback)
	app.Handle("GET", "/logout/:provider", LogoutHandler)
	app.Handle("GET", "/register", RegisterGet)
	app.Handle("POST", "/register", db.RegisterPost)

	requiresLogin.Handle("GET", "/dashboard", db.DashboardGet)
	requiresLogin.Handle("GET", "/dashboard/other", db.DashboardOtherGet)

	requiresAdmin.Handle("GET", "/admin", db.AdminGet)
	requiresAdmin.Handle("GET", "/admin/reviewed", db.AdminReviewedGet)
	requiresAdmin.Handle("GET", "/admin/edit/:EID", db.EditGet)
	requiresAdmin.Handle("POST", "/admin/edit/:EID", db.EditPost)
	app.Handle("GET", "/fetchall", db.FetchLatest)
	requiresAdmin.Handle("GET", "/refetch/:Eid", db.RefetchEid)

	return app
}
