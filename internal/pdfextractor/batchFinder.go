package pdfextractor

import (
	"regexp"
)

type Location struct {
	loc int
	val string
}

func FindBatch(documentBody []string, year string) string {
	locs := findBatchLocations(documentBody)
	nbrs := findNeighbours(locs, documentBody)
	val := findBatchYear(year, nbrs)
	return val
}

// this function tries to find the location of word 'batch' in the pdf
func findBatchLocations(body []string) []Location {
	var count int
	locations := make([]Location, 0, 10)
	r := regexp.MustCompile(`^(?i)(batch)?$`)
	for k, v := range body {
		// improve by adding regex for the word batch here and
		// then that regex will be compared for each word v in the body
		if r.MatchString(v) {
			count = count + 1
			location := Location{
				loc: k,
				val: v,
			}
			locations = append(locations, location)
		}
	}
	//fmt.Println("Batch Locations are: ", locations)

	// fmt.Println("Count is", count)
	// fmt.Println("Location are, ", locations)
	return locations
}

type Neighbour struct {
	n string
}

// this function find the neighbours of the word batch so that we can find the number
func findNeighbours(loc []Location, s []string) []Neighbour {
	neighbours := make([]Neighbour, 0, 14)
	for i := 0; i < len(loc); i++ {
		pos := loc[i].loc

		for j := pos - 3; j <= pos+3; j++ {
			wl := findWordAtPosition(j, s)
			neighbour := Neighbour{
				n: wl,
			}
			neighbours = append(neighbours, neighbour)
		}

	}
	//fmt.Println("NEIGHBOURS ARE: ", neighbours)
	return neighbours
}

func findWordAtPosition(i int, s []string) string {
	for k, v := range s {
		if k == i {
			return v
		}
	}
	return ""
}

// this function find the batch year given in those neighbours
func findBatchYear(year string, nbrs []Neighbour) string {
	var count int
	len := len(nbrs)

	for i := 0; i < len; i++ {
		if nbrs[i].n == year {
			count = count + 1
		}
	}

	// fmt.Println("Counts of 2018 are: ", count)
	if count > 0 {
		return year
	}
	return ""
}
