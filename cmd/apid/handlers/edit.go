package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/amankapoor/placementpal/internal/pdfextractor"
	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/amankapoor/placementpal/internal/platform/web"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (db *Database) EditGet(ctx context.Context, w http.ResponseWriter, r *http.Request, ps web.Params) error {
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

	eid := ps.ByName("EID")

	pd, err := fetchPlacementByEID(ctx, reqDB, eid)
	if err != nil {
		return errors.Wrapf(err, "Could not fetch btech data")
	}
	ud, err := fetchUserData(ctx, reqDB, gID)
	if err != nil {
		return errors.Wrapf(err, "Could not fetch user data")
	}
	var dgs []string
	for _, v := range pdfextractor.Degrees {
		s := strings.TrimPrefix(v.Name, "{")
		s = strings.TrimPrefix(v.Name, "}")
		//fmt.Println(k, s)
		dgs = append(dgs, s)
	}
	data := map[string]interface{}{
		"placement": pd,
		"user":      &ud,
		"degrees":   dgs,
	}
	web.Respond(ctx, w, "adminEdit.html", data, http.StatusOK)
	return nil
}

type updates struct {
	EID       string `valid:"required"`
	Batch     []string
	Degrees   []string
	CTC       string
	DriveDate time.Time
}

func (db *Database) EditPost(ctx context.Context, w http.ResponseWriter, r *http.Request, _ web.Params) error {
	err := r.ParseForm()
	if err != nil {
		return errors.Wrap(err, "unable to parse register form")
	}

	eid := r.FormValue("eid")
	// fmt.Println(eid)
	batch := r.Form["batch"]
	// fmt.Println(batch)

	degrees := r.Form["degrees"]
	// fmt.Println(degrees)

	ctc := r.FormValue("ctc")
	// fmt.Println(ctc)

	layout := "2006-01-02"
	driveDate := r.FormValue("driveDate")
	// fmt.Println(driveDate)
	var t time.Time
	if driveDate == "" {
		t = time.Time{}
	} else {
		t, err = time.Parse(layout, driveDate)
		if err != nil {
			return errors.Wrap(err, "unable to parse drive date to different time")
		}
	}

	// fmt.Println(t)

	s := updates{
		EID:       eid,
		Batch:     batch,
		Degrees:   degrees,
		CTC:       ctc,
		DriveDate: t,
	}

	//fmt.Println("S is: ", s)

	res, err := govalidator.ValidateStruct(s)
	if err != nil || res == false {
		return errors.Wrap(err, "invalid inputs by user")
	}
	//fmt.Println(res)

	reqDB, err := db.MasterDB.MGOCopy()
	if err != nil {
		return errors.Wrapf(web.ErrDBNotConfigured, "")
	}
	defer reqDB.MGOClose()

	err = update(ctx, reqDB, s)
	if err != nil {
		return errors.Wrapf(err, "<<Unable to update placement eid %s>>", eid)
	}

	w.Write([]byte("ok"))
	return nil
}

func update(ctx context.Context, dbConn *db.DB, s updates) error {
	url, err := strconv.Atoi(s.EID)
	if err != nil {
		return errors.Wrap(err, "<<Unable to convert Atoi while updating placement>>")
	}

	//fmt.Println("S DRIVE DATE ", s.DriveDate)
	q := bson.M{"url": url}
	m := bson.M{
		"$set": bson.M{
			"degrees":     s.Degrees,
			"batch":       s.Batch,
			"ctc":         s.CTC,
			"drive_date":  s.DriveDate,
			"is_reviewed": true,
		},
	}

	f := func(collection *mgo.Collection) error {
		return collection.Update(q, m)
	}
	if err := dbConn.MGOExecute(ctx, masterDataCollection, f); err != nil {
		if err == mgo.ErrNotFound {
			return web.ErrNotFound
		}
		return errors.Wrap(err, fmt.Sprintf("db.users.update(%s, %s)", db.Query(q), db.Query(m)))
	}

	return nil
}

func fetchPlacementByEID(ctx context.Context, dbConn *db.DB, eid string) (*DisplayPlacementData, error) {
	var md *DisplayPlacementData
	url, err := strconv.Atoi(eid)
	if err != nil {
		return nil, err
	}
	query := bson.M{
		"url": url,
	}
	f := func(collection *mgo.Collection) error {
		return collection.Find(query).One(&md)
	}

	if err := dbConn.MGOExecute(ctx, masterDataCollection, f); err != nil {
		return nil, errors.Wrap(err, "unable to get one placement data for admin edit")
	}

	return md, nil
}
