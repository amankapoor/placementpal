package handlers

import (
	"context"
	"net/http"

	"github.com/amankapoor/placementpal/internal/platform/web"
)

func Home(ctx context.Context, w http.ResponseWriter, r *http.Request, _ web.Params) error {
	web.Respond(ctx, w, "home.html", ProviderIndexVar, http.StatusOK)
	return nil
}
