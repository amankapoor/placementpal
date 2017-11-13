package handlers

import (
	"context"
	"net/http"
	"os"

	"github.com/amankapoor/placementpal/internal/files"
	"github.com/amankapoor/placementpal/internal/pdfextractor"
	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/amankapoor/placementpal/internal/platform/web"
	"github.com/amankapoor/placementpal/internal/scraper"
	"github.com/pkg/errors"
)

type Database struct {
	MasterDB *db.DB
}

// FetchLatest fetches the latest entries from placement site
// and adds the beautified-retrieved data to latest_fetch collection
func (db *Database) FetchLatest(ctx context.Context, w http.ResponseWriter, r *http.Request, _ web.Params) error {

	//fetching beautified data
	beautifiedData := scraper.FetchLatest()
	//fmt.Printf("\n\n\n BEAUTIFIED DATA IS: %v\n\n\n", beautifiedData)

	reqDB, err := db.MasterDB.MGOCopy()
	if err != nil {
		return errors.Wrapf(web.ErrDBNotConfigured, "")
	}
	defer reqDB.MGOClose()

	//fmt.Println("after mgo close")
	// adding beautified data to latest_fetch collection
	err = scraper.AddLatestFetchToDB(ctx, reqDB, beautifiedData)
	if err != nil {
		return errors.Wrapf(err, "Could not add latest fetch to db")
	}

	//fmt.Println("after doing AddLatestFetchToDB call")
	//web.Respond(ctx, w, "dashboard.html", nil, http.StatusOK)

	// appnding uniques(new) to master data
	uniques, err := scraper.AppendUniquesToMaster(ctx, reqDB)
	if err != nil {
		return errors.Wrapf(err, "<<Could not append uniques from latest_fetch to master_data>>")
	}

	// updating the changed ones in master data
	existing, err := scraper.UpdateExistingIfAnyInMasterData(ctx, reqDB)
	if err != nil {
		return errors.Wrapf(err, "<<Unable to update existings in master data>>")
	}

	// These Eids include the new and the updated ones for downloading
	var Eids []int
	for _, v := range uniques {
		Eids = append(Eids, v)
	}
	for _, v := range existing {
		Eids = append(Eids, v.URL)
	}
	//fmt.Printf("\n\nEIDS ARE: %v", Eids)
	var fnames []string
	var fileNames []string
	for _, v := range Eids {
		// downloads files to temporary path
		fileNames, err = files.Downloader(v)
		if err != nil {
			return errors.Wrap(err, "<<unable to download the new Eids>>")
		}
		var degrees []string
		var batches []string
		var btchs []string
		var dgrs []string
		for _, v := range fileNames {
			btchs, dgrs, err = pdfextractor.DataFromPDF(v)
			if err != nil {
				return errors.Wrap(err, "<<Unable to extract data from file>>")
			}

			batches = append(batches, btchs...)
			degrees = append(degrees, dgrs...)
			fnames = append(fnames, v)
		}

		//fmt.Println("DEGREES AND BATCHES ARE: ", degrees, batches)
		var tempPath string
		err = files.AddFileInfoToMasterData(ctx, reqDB, v, fileNames, batches, degrees)
		if err != nil {
			// here remove the files from temp path first
			for _, v := range fileNames {
				tempPath = "/home/amankapoor/go/src/github.com/amankapoor/placementpal/cmd/apid/views/temp/" + v
				err := os.Remove(tempPath)
				if err != nil {
					return errors.Wrapf(err, "<<unable to delete filename: %s from temp directory. Delete it manually.>>", v)
				}
			}
			// and then return the database error
			return errors.Wrapf(err, "<<unable to add file name %s to master data>>", fileNames)
		}
	}

	for _, v := range fnames {
		// now since data is added successfully to master data so move files to files path
		tempPath := "/home/amankapoor/go/src/github.com/amankapoor/placementpal/cmd/apid/views/temp/" + v
		newPath := "/home/amankapoor/go/src/github.com/amankapoor/placementpal/cmd/apid/files/pdf/" + v
		err = os.Rename(tempPath, newPath)
		if err != nil {
			return errors.Wrapf(err, "<<unable to move filename: %s from temp directory to files directory. Returning error.>>", v)
		}
	}

	// err, info := files.Downloader(listOfNewEntries)
	// if err != nil {
	// 	return errors.Wrapf(err, "Unable to download all the files of all the new entries")
	// }
	w.Write([]byte("ok"))
	return nil
}
