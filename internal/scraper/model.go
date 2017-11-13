package scraper

import "time"

type placementData struct {
	title string
	url   string
}

// BeautifiedPlacementData contains sorted (by ints)
// and clean info of all latest fetched placements
type BeautifiedPlacementData struct {
	Title string
	URL   int
}

// CreateLatestFetch contains information about latest fetches to feed to database
type CreateLatestFetch struct {
	Title string
	URL   int
}

type CreateExistingData = CreateLatestFetch

type LatestFetch struct {
	ID          string     `bson:"_id"`
	Title       string     `bson:"title"`
	URL         int        `bson:"url"`
	DateCreated *time.Time `bson:"date_created"`
}

type MasterData = LatestFetch
type Unique = LatestFetch
type Existing = LatestFetch
