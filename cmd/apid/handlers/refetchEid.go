package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/amankapoor/placementpal/internal/platform/web"
	"github.com/amankapoor/placementpal/internal/scraper"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (db *Database) RefetchEid(ctx context.Context, w http.ResponseWriter, r *http.Request, ps web.Params) error {

	eid := ps.ByName("Eid")

	reqDB, err := db.MasterDB.MGOCopy()
	if err != nil {
		return errors.Wrapf(web.ErrDBNotConfigured, "")
	}
	defer reqDB.MGOClose()

	intEid, newTitleToUpdate, newFileNamesToUpdate, newFileNamesForDeletionInCaseOpFails, newDegreesToUpdate, newBatchToUpdate, err := scraper.RefetchSingleEid(ctx, reqDB, eid)
	if err != nil {
		return errors.Wrapf(err, "<<Unable to refetch single Eid: %s>>", eid)
	}

	err = addSingleEIDFileInfoToMasterData(ctx, reqDB, intEid, newFileNamesToUpdate, newBatchToUpdate, newDegreesToUpdate, newTitleToUpdate)
	if err != nil {
		// here remove the files from temp path first
		for _, v := range newFileNamesForDeletionInCaseOpFails {
			tempPath := "/home/amankapoor/go/src/github.com/amankapoor/placementpal/cmd/apid/views/temp/" + v
			err := os.Remove(tempPath)
			if err != nil {
				return errors.Wrapf(err, "<<unable to delete filename: %s from temp directory. Delete it manually.>>", v)
			}
		}
		// and then return the database error
		return errors.Wrap(err, "<<unable to add single eid info to master data>>")
	}

	// this newFileNamesForDeletionInCaseOpFails is also the new files that
	// are downloaded
	for _, v := range newFileNamesForDeletionInCaseOpFails {
		// now since data is added successfully to master data so move files to files path
		tempPath := "/home/amankapoor/go/src/github.com/amankapoor/placementpal/cmd/apid/views/temp/" + v
		newPath := "/home/amankapoor/go/src/github.com/amankapoor/placementpal/cmd/apid/views/pdf/" + v
		err = os.Rename(tempPath, newPath)
		if err != nil {
			return errors.Wrapf(err, "<<unable to move filename: %s from temp directory to files directory. Returning error.>>", v)
		}
	}

	w.Write([]byte("ok single refetch"))
	return nil
}

const masterDataCollection = "master_data"

func addSingleEIDFileInfoToMasterData(ctx context.Context, dbConn *db.DB, Eid int, fileNames []string, batches []string, degrees []string, title string) error {
	now := time.Now()

	selector := bson.M{"url": Eid}
	query := bson.M{
		"title":         title,
		"files":         fileNames,
		"batch":         batches,
		"degrees":       degrees,
		"date_modified": &now,
	}
	update := bson.M{"$set": query}

	f := func(collection *mgo.Collection) error {
		return collection.Update(selector, update)
	}

	if err := dbConn.MGOExecute(ctx, masterDataCollection, f); err != nil {
		return errors.Wrap(err, fmt.Sprintf("db.users.update(%s)", db.Query(update)))
	}

	return nil
}
