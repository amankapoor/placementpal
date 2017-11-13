package handlers

import (
	"context"
	"net/http"

	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/amankapoor/placementpal/internal/platform/web"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
)

func (db *Database) AdminGet(ctx context.Context, w http.ResponseWriter, r *http.Request, _ web.Params) error {
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

	// if gID != "116042618152348131594" {
	// 	http.Error(w, "This page does not exist.", 404)
	// 	return nil
	// }

	pd, err := fetchAllPlacements(ctx, reqDB)
	if err != nil {
		return errors.Wrapf(err, "Could not fetch btech data")
	}

	ud, err := fetchUserData(ctx, reqDB, gID)
	if err != nil {
		return errors.Wrapf(err, "Could not fetch user data")
	}

	data := map[string]interface{}{
		"placement": pd,
		"user":      &ud,
	}

	//fmt.Println(data["placement"], data["user"])
	web.Respond(ctx, w, "admin.html", data, http.StatusOK)
	return nil
}

func fetchAllPlacements(ctx context.Context, dbConn *db.DB) ([]*DisplayPlacementData, error) {
	var md []*DisplayPlacementData

	f := func(collection *mgo.Collection) error {
		return collection.Find(nil).Sort("-url").All(&md)
	}

	if err := dbConn.MGOExecute(ctx, masterDataCollection, f); err != nil {
		return nil, errors.Wrap(err, "unable to get full master data for admin")
	}

	return md, nil
}
