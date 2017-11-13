package scraper

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CheckExistingForChanges finds the existing ones which are changed
// and updates them in the master data
func UpdateExistingIfAnyInMasterData(ctx context.Context, dbConn *db.DB) ([]CreateExistingData, error) {
	// get master data
	mdurls, err := getMasterDataURLs(ctx, dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get master data urls")
	}
	//fmt.Println("Printing masterdata", mdurls)

	// get latest data
	ldurls, err := getLatestDataURLs(ctx, dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get latest data urls")
	}
	//fmt.Println("Printing latestdata", ldurls)

	// findExisting finds existing which are in master_data collection
	existing := findExisting(mdurls, ldurls)
	//fmt.Println("Printing existing data", existing)

	// get full latest_fetch
	flf, err := getFullLatestFetch(ctx, dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get latest fetched data")
	}

	// get full master data
	fmd, err := getFullMasterData(ctx, dbConn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get master data")
	}
	//fmt.Println("Printing full master data", fmd)

	// adding uniques to master data
	updated, err := updateExistingInMasterData(ctx, dbConn, existing, flf, fmd)
	if err != nil {
		return nil, errors.Wrap(err, "unable to update existing in master data")
	}

	return updated, nil
}

// findExisting returns the exisitng values
// which are both in master and latest data
func findExisting(masterData []int, latestData []int) []int {
	var diff []int
	var v int
	for _, v = range latestData {
		i := sort.SearchInts(masterData, v)
		if i < len(masterData) && masterData[i] == v {
			diff = append(diff, v)
			//fmt.Printf("\nvalue v %d exists in master data at position i %d", v, i)
		} else {
			// fmt.Printf("\nvalue v %d does not exist in data and would be inserted at %d", v, i)
			// diff = append(diff, v)
			continue
		}
	}
	return diff
}

// updateExistingInMasterData updates the existing documents in master_data
// to reflect latest changes
func updateExistingInMasterData(ctx context.Context, dbConn *db.DB, existingInts []int, latestData []LatestFetch, masterData []MasterData) ([]CreateExistingData, error) {
	existings := rebuildExisting(existingInts, latestData)
	//fmt.Println("After rebuilding Existing are: ", existings)

	newData := compareTitleStringsOfExistingsAndMaster(existings, masterData)
	//fmt.Println("New data to update is: ", newData)

	now := time.Now()
	for _, v := range newData {
		selector := bson.M{"url": v.URL}
		query := bson.M{
			"title":         v.Title,
			"url":           v.URL,
			"is_reviewed":   false,
			"date_modified": &now,
		}
		update := bson.M{"$set": query}

		f := func(collection *mgo.Collection) error {
			return collection.Update(selector, update)
		}

		if err := dbConn.MGOExecute(ctx, masterDataCollection, f); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("db.users.update(%s)", db.Query(update)))
		}
	}
	return newData, nil
}

// compareTitleStringsOfExistingsAndMaster compares title strings
// of existing and master data and returns only the changes ones
func compareTitleStringsOfExistingsAndMaster(existing []Existing, masterData []MasterData) []CreateExistingData {
	var ced []CreateExistingData
	for _, v := range existing {
		for _, iv := range masterData {
			if v.URL == iv.URL {
				res := strings.Compare(v.Title, iv.Title)
				if res != 0 {
					//fmt.Printf("in the loop comparing %v and %v", v.Title, iv.Title)
					a := CreateExistingData{
						Title: v.Title,
						URL:   v.URL,
					}
					ced = append(ced, a)
					// fmt.Println("result is: ", res)
					// fmt.Println("appended a ", a)
					// fmt.Println("")
				}
			}
		}
	}
	return ced
}

// rebuildExisting is needed because in earlier steps
// we play with only int urls, but to add in master_data
// we need all the fields. So, it ranges over latest fetched data
// and compares if its url is unique to us, if it finds any unique url
// it appends all the data of that url to new list that it outputs for inclusion
// in master_data collection
func rebuildExisting(ui []int, ld []LatestFetch) []Existing {
	var rebuilt []Existing

	for _, v := range ui {
		for _, iv := range ld {
			if v == iv.URL {
				rebuilt = append(rebuilt, iv)
			}
		}
	}
	// now := time.Now()

	// new2 := Existing{
	// 	DateCreated: &now,
	// 	ID:          "2",
	// 	URL:         2,
	// 	Title:       "SANSAR2",
	// }

	// rebuilt = append(rebuilt, new2)
	return rebuilt
}
