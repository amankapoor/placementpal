package scraper

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const masterDataCollection = "master_data"

// AppendMaster computes the difference between latest_fetch database
// and master-data, and updates the master_data only with new entries
// i.e., the entries which do not exist in master db.
//
// Then, in second phase, it checks for differences in the titles of master_data
// and if discovers any changes, adds new version below to the document
func AppendUniquesToMaster(ctx context.Context, dbConn *db.DB) ([]int, error) {

	// get master data
	mdurls, err := getMasterDataURLs(ctx, dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "<<unable to get master data urls>>")
	}
	//fmt.Println("Printing masterdata", mdurls)

	// get latest data
	ldurls, err := getLatestDataURLs(ctx, dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "<<unable to get latest data urls>>")
	}
	//fmt.Println("Printing latestdata", ldurls)

	// findUniques finds uniques which are not in master_data collection
	uniques := findUniques(mdurls, ldurls)
	//fmt.Println("Printing unique data", uniques)

	// get full latest_fetch
	flf, err := getFullLatestFetch(ctx, dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "<<unable to get latest fetched data>>")
	}

	// adding uniques to master data
	err = addUniquesToMasterData(ctx, dbConn, uniques, flf)
	if err != nil {
		return nil, errors.Wrap(err, "<<unable to add uniques to master data>>")
	}

	return uniques, nil

	// // get full master data to test
	// fmd, err := getFullMasterData(ctx, dbConn)
	// if err != nil {
	// 	return errors.Wrap(err, "unable to get master data")
	// }
	// fmt.Println("Printing masterdata again", fmd)
}

// getMasterData gets urls ints from master_data collection
func getMasterDataURLs(ctx context.Context, dbConn *db.DB) ([]int, error) {
	query := bson.M{"url": 1}
	var masterData []int

	var md []struct {
		URL int `bson:"url"`
	}

	f := func(collection *mgo.Collection) error {
		return collection.Find(nil).Select(query).Sort("url").All(&md)
	}

	if err := dbConn.MGOExecute(ctx, masterDataCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.master_data.find({}).select(%s).all()", db.Query(query)))
	}

	for _, v := range md {
		masterData = append(masterData, v.URL)
	}

	return masterData, nil
}

// getLatestData gets urls ints of latest_fetch collection
func getLatestDataURLs(ctx context.Context, dbConn *db.DB) ([]int, error) {
	query := bson.M{"url": 1}
	var latestData []int

	var ld []struct {
		URL int `bson:"url"`
	}

	f := func(collection *mgo.Collection) error {
		return collection.Find(nil).Select(query).Sort("url").All(&ld)
	}

	if err := dbConn.MGOExecute(ctx, latestFetchCollection, f); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("db.master_data.find({}).select(%s).all()", db.Query(query)))
	}

	for _, v := range ld {
		latestData = append(latestData, v.URL)
	}

	return latestData, nil
}

// findUniques returns the unique values
// which are IN latest data but NOT IN master data
func findUniques(masterData []int, latestData []int) []int {
	var diff []int
	var v int
	for _, v = range latestData {
		i := sort.SearchInts(masterData, v)
		if i < len(masterData) && masterData[i] == v {
			//fmt.Printf("\nvalue v %d exists in master data at position i %d", v, i)
			continue
		} else {
			//fmt.Printf("\nvalue v %d does not exist in data and would be inserted at %d", v, i)
			diff = append(diff, v)
		}
	}
	return diff
}

// getFullLatestFetch gets full latest_fetch collection
func getFullLatestFetch(ctx context.Context, dbConn *db.DB) ([]LatestFetch, error) {

	var lf []LatestFetch

	f := func(collection *mgo.Collection) error {
		return collection.Find(nil).All(&lf)
	}

	if err := dbConn.MGOExecute(ctx, latestFetchCollection, f); err != nil {
		return nil, errors.Wrap(err, "unable to get full latest_fetch collection")
	}

	return lf, nil
}

// getFullMasterData gets full master_data collection
func getFullMasterData(ctx context.Context, dbConn *db.DB) ([]MasterData, error) {

	var lf []MasterData

	f := func(collection *mgo.Collection) error {
		return collection.Find(nil).All(&lf)
	}

	if err := dbConn.MGOExecute(ctx, masterDataCollection, f); err != nil {
		return nil, errors.Wrap(err, "unable to get full master_data collection")
	}

	return lf, nil
}

// addUniquesToMasterData inserts the unique documents to master_data
func addUniquesToMasterData(ctx context.Context, dbConn *db.DB, uniqueInts []int, latestData []LatestFetch) error {
	uniques := rebuildUniques(uniqueInts, latestData)
	//fmt.Println("Appending these uniques to master ", uniques)

	now := time.Now()
	for _, v := range uniques {
		query := bson.M{
			"_id":           v.ID,
			"title":         v.Title,
			"url":           v.URL,
			"is_reviewed":   false,
			"date_created":  &now,
			"date_modified": &now,
		}

		f := func(collection *mgo.Collection) error {
			return collection.Insert(query)
		}

		if err := dbConn.MGOExecute(ctx, masterDataCollection, f); err != nil {
			return errors.Wrap(err, fmt.Sprintf("db.users.insert(%s)", db.Query(query)))
		}
	}
	return nil
}

// rebuildUniques is needed because in earlier steps
// we play with only int urls, but to add in master_data
// we need all the fields. So, it ranges over latest fetched data
// and compares if its url is unique to us, if it finds any unique url
// it appends all the data of that url to new list that it outputs for inclusion
// in master_data collection
func rebuildUniques(ui []int, ld []LatestFetch) []Unique {
	var rebuilt []Unique

	for _, v := range ui {
		for _, iv := range ld {
			if v == iv.URL {
				rebuilt = append(rebuilt, iv)
			}
		}
	}
	// now := time.Now()

	// new2 := Unique{
	// 	DateCreated: &now,
	// 	ID:          "2",
	// 	URL:         2,
	// 	Title:       "SANSAR2",
	// }

	// rebuilt = append(rebuilt, new2)
	return rebuilt
}
