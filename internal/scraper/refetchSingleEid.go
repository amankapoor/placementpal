package scraper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/amankapoor/placementpal/internal/files"
	"github.com/amankapoor/placementpal/internal/pdfextractor"
	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MasterDataForSingleEIDRefetch struct {
	ID           string    `bson:"_id"`
	URL          int       `bson:"url"`
	DateCreated  time.Time `bson:"date_created"`
	DateModified time.Time `bson:"date_modified"`
	Title        string    `bson:"title"`
	Files        []string  `bson:"files"`
	Batch        []string  `bson:"batch"`
	Degrees      []string  `bson:"degrees"`
}

func RefetchSingleEid(ctx context.Context, dbConn *db.DB, eid string) (int, string, []string, []string, []string, []string, error) {

	// fetchlatest to find title of the eid there
	beautifiedData := FetchLatest()

	// find the title in the data
	var newTitle string
	urlID, err := strconv.Atoi(eid)
	if err != nil {
		return 0, "", nil, nil, nil, nil, errors.Wrapf(err, "<<unable to convert ascii eid: %s to int>>", eid)
	}
	for _, v := range beautifiedData {
		if v.URL == urlID {
			newTitle = v.Title
		}
	}
	if newTitle == "" {
		return 0, "", nil, nil, nil, nil, errors.Wrapf(err, "<<title not found while refetching single Eid: %s>>", eid)
	}

	// download its files
	newFileNames, err := files.Downloader(urlID)
	if err != nil {
		return 0, "", nil, nil, nil, nil, errors.Wrap(err, "<<unable to download files of the new Eid>>")
	}

	// extract new data from its files
	var degrees []string
	var batches []string
	var btchs []string
	var dgrs []string
	for _, v := range newFileNames {
		btchs, dgrs, err = pdfextractor.DataFromPDF(v)
		if err != nil {
			return 0, "", nil, nil, nil, nil, errors.Wrap(err, "<<Unable to extract data from file>>")
		}

		batches = append(batches, btchs...)
		degrees = append(degrees, dgrs...)
	}

	// get current master data about this eid
	q := bson.M{"url": urlID}
	var u MasterDataForSingleEIDRefetch
	f := func(collection *mgo.Collection) error {
		return collection.Find(q).One(&u)
	}
	if err := dbConn.MGOExecute(ctx, masterDataCollection, f); err != nil {
		return 0, "", nil, nil, nil, nil, errors.Wrap(err, fmt.Sprintf("Entry not found in DB - db.users.find(%s)", db.Query(q)))
	}

	// preparing variables to update the database
	//var n RefetchSingleEID
	//fmt.Println("OLD DATA IS: ", &u)

	newFileNamesToUpdate, newDegreesToUpdate, newBatchToUpdate := appropriateFileNames(newFileNames, u.Files, u.Degrees, u.Batch, degrees, batches)
	// next we change the title to new title
	newTitleToUpdate := newTitle

	// return these these new variables to add t master data
	return urlID, newTitleToUpdate, newFileNamesToUpdate, newFileNames, newDegreesToUpdate, newBatchToUpdate, nil
}

// if there are no new files then we return old,
// else we append new filenames to previous filenames
func appropriateFileNames(newDownloadedFileNames, oldFileNames, oldDegrees, oldBatch, newDegrees, newBatches []string) ([]string, []string, []string) {
	var n []string
	if len(newDownloadedFileNames) == 0 {
		return oldFileNames, oldDegrees, oldBatch
	}

	n = append(newDownloadedFileNames, oldFileNames...)
	//fmt.Println("NEW FILE NAMES ARE: ", n)
	return n, newDegrees, newBatches
}
