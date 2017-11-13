package scraper

import (
	"context"
	"fmt"
	"time"

	"github.com/amankapoor/placementpal/internal/platform/db"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const latestFetchCollection = "latest_fetch"

// AddLatestFetchToDB adds the latest fetched data to latest_fetch collection
func AddLatestFetchToDB(ctx context.Context, dbConn *db.DB, clf []BeautifiedPlacementData) error {

	err := emptyCollection(ctx, latestFetchCollection, dbConn)
	if err != nil {
		return errors.Wrap(err, "cannot remove existing documents from the collection")
	}

	now := time.Now()
	//fmt.Println("after time now 3")
	var query bson.M
	for _, v := range clf {
		// fmt.Println("doing for range for k ", k)
		query = bson.M{
			"_id":          bson.NewObjectId().Hex(),
			"title":        v.Title,
			"url":          v.URL,
			"date_created": &now,
		}
		//fmt.Println("INSERTING TO LATEST_FETCH: ", k)
		f := func(collection *mgo.Collection) error {
			return collection.Insert(query)
		}

		if err := dbConn.MGOExecute(ctx, latestFetchCollection, f); err != nil {
			return errors.Wrap(err, fmt.Sprintf("db.users.insert(%s)", db.Query(query)))
		}
	}

	return nil
}

// emptyCollection removes all the documents from the said (latest_fetch) collection
func emptyCollection(ctx context.Context, collectioName string, dbConn *db.DB) error {
	query := bson.M{}
	f := func(collection *mgo.Collection) error {
		_, err := collection.RemoveAll(query)
		if err != nil {
			return errors.Wrap(err, "cannot remove existing documents from the collection")
		}
		return nil
	}

	if err := dbConn.MGOExecute(ctx, latestFetchCollection, f); err != nil {
		return errors.Wrap(err, fmt.Sprintf("db.users.remove(%s)", db.Query(query)))
	}
	return nil
}
