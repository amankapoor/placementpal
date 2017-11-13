package handlers

import (
	"context"
	"net/http"

	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/amankapoor/placementpal/internal/platform/web"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (db *Database) DashboardOtherGet(ctx context.Context, w http.ResponseWriter, r *http.Request, _ web.Params) error {
	//fmt.Println(r.Method)
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
	batch, err := s.GetString(SVBatch)
	if err != nil {
		return errors.Wrapf(err, "Could not fetch batch session")
	}
	// degree, err := s.GetString(SVDegree)
	// if err != nil {
	// 	return errors.Wrapf(err, "Could not fetch degree from session")
	// }
	s.PopString(w, SVFirstName)

	pd, err := fetchOtherPlacements(ctx, reqDB, batch)
	if err != nil {
		return errors.Wrapf(err, "Could not fetch btech data")
	}

	ud, err := fetchUserData(ctx, reqDB, gID)
	if err != nil {
		return errors.Wrapf(err, "Could not fetch user data")
	}
	// //fmt.Println(md)
	// s := sessionManager.Load(r)
	// fname, err := s.GetString(SVFirstName)
	// if err != nil {
	// 	return errors.Wrapf(err, "Could not fetch first name data")
	// }
	// profilePhoto, err := s.GetString(SVProfilePhoto)
	// if err != nil {
	// 	return errors.Wrapf(err, "Could not fetch profile photo from session")
	// }

	placement := make([]PlacementDataToDisplay, len(pd))

	for k, v := range pd {
		placement[k].CTC = v.CTC
		placement[k].Files = v.Files
		placement[k].Title = v.Title
		// PlacementDataToDisplay[k].Title = pd[k].Title
		// PlacementDataToDisplay[k].Files = pd[k].Files
		// PlacementDataToDisplay[k].CTC = pd[k].CTC
		var s string
		if !v.DriveDate.IsZero() {
			s = v.DriveDate.Format("02-Jan-2006")
			// PlacementDataToDisplay[k].DriveDate = s

			placement[k].DriveDate = s
		}

	}

	//fmt.Println(placement)

	data := map[string]interface{}{
		"placement": placement,
		"user":      &ud,
	}

	//fmt.Println(data["placement"], data["user"])
	web.Respond(ctx, w, "dashboardOther.html", data, http.StatusOK)
	return nil
}

func fetchOtherPlacements(ctx context.Context, dbConn *db.DB, batch string) ([]*DisplayPlacementData, error) {
	var md []*DisplayPlacementData

	query := bson.M{
		"$and": []bson.M{
			bson.M{
				"degrees": []string{},
			},
			bson.M{
				"batch": bson.M{
					"$in": []string{batch},
				},
			},
		},
	}

	f := func(collection *mgo.Collection) error {
		return collection.Find(query).Limit(20).Sort("-url").All(&md)
	}

	if err := dbConn.MGOExecute(ctx, masterDataCollection, f); err != nil {
		return nil, errors.Wrap(err, "unable to get find and of degree and batch")
	}

	return md, nil
}
