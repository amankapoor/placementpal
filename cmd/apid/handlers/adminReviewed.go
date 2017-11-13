package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/amankapoor/placementpal/internal/platform/web"
	"github.com/pkg/errors"
)

type AdminReviewedPlacementDataToDisplay struct {
	ID           string    `bson:"_id"`
	URL          int       `bson:"url"`
	Title        string    `bson:"title"`
	Files        []string  `bson:"files"`
	Degrees      []string  `bson:"degrees"`
	Batches      []string  `bson:"batch"`
	DateModified time.Time `bson:"date_modified"`
	CTC          string    `bson:"ctc"`
	DriveDate    string    `bson:"drive_date"`
	IsReviewed   bool      `bson:"is_reviewed"`
}

func (db *Database) AdminReviewedGet(ctx context.Context, w http.ResponseWriter, r *http.Request, _ web.Params) error {
	reqDB, err := db.MasterDB.MGOCopy()
	if err != nil {
		return errors.Wrapf(web.ErrDBNotConfigured, "")
	}
	defer reqDB.MGOClose()

	s := sessionManager.Load(r)
	gID, err := s.GetString(SVGoogleID)
	if err != nil {
		return errors.Wrapf(err, "Could not fetch google ID from session")
	}

	if gID != "116042618152348131594" {
		return errors.Wrap(nil, "<<somebody else trying to access admin panel>>")
	}

	pd, err := fetchAllPlacements(ctx, reqDB)
	if err != nil {
		return errors.Wrapf(err, "Could not fetch btech data")
	}

	ud, err := fetchUserData(ctx, reqDB, gID)
	if err != nil {
		return errors.Wrapf(err, "Could not fetch user data")
	}

	pm := make([]AdminReviewedPlacementDataToDisplay, len(pd))

	for k, v := range pd {
		pm[k].CTC = v.CTC
		pm[k].Files = v.Files
		pm[k].Title = v.Title
		pm[k].Batches = v.Batches
		pm[k].Degrees = v.Degrees
		pm[k].IsReviewed = v.IsReviewed
		pm[k].URL = v.URL
		var s string
		if !v.DriveDate.IsZero() {
			s = v.DriveDate.Format("02-Jan-2006")
			// PlacementDataToDisplay[k].DriveDate = s

			pm[k].DriveDate = s
		}
	}

	//fmt.Println("PM is: ", pm)

	data := map[string]interface{}{
		"placement": pm,
		"user":      &ud,
	}

	//fmt.Println(data["placement"], data["user"])
	web.Respond(ctx, w, "adminReviewed.html", data, http.StatusOK)
	return nil
}
