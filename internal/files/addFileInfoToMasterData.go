package files

import (
	"context"
	"fmt"
	"time"

	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const masterDataCollection = "master_data"

func AddFileInfoToMasterData(ctx context.Context, dbConn *db.DB, Eid int, fileNames []string, batches []string, degrees []string) error {
	now := time.Now()

	selector := bson.M{"url": Eid}
	query := bson.M{
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
