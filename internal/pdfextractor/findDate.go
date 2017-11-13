package pdfextractor

// we are not using find date until we improve the code
// import (
// 	"fmt"
// 	"regexp"
// 	"strconv"
// 	"time"
// )

// func FindDate(documentBody []string) string {
// 	locs := findBatchLocations(documentBody)
// 	nbrs := findDateNeighbours(locs, documentBody)
// 	year := strconv.Itoa(time.Now().Year())
// 	val := findDateInNeighbours(year, nbrs)
// 	return val

// }

// // this function tries to find the location of word 'date' in the pdf
// func findDateLocations(body []string) []Location {
// 	var count int
// 	locations := make([]Location, 0, 10)
// 	r := regexp.MustCompile(`^(?i)(date)?$`)
// 	for k, v := range body {
// 		// improve by adding regex for the word batch here and
// 		// then that regex will be compared for each word v in the body
// 		if r.MatchString(v) {
// 			count = count + 1
// 			location := Location{
// 				loc: k,
// 				val: v,
// 			}
// 			locations = append(locations, location)
// 		}
// 	}
// 	//fmt.Println("Date Locations are: ", locations)

// 	// fmt.Println("Count is", count)
// 	// fmt.Println("Location are, ", locations)
// 	return locations
// }

// // this function find the neighbours of the word date so that we can find the number
// func findDateNeighbours(loc []Location, s []string) []Neighbour {
// 	neighbours := make([]Neighbour, 0, 14)
// 	for i := 0; i < len(loc); i++ {
// 		pos := loc[i].loc

// 		for j := pos; j <= pos+10; j++ {
// 			wl := findWordAtPosition(j, s)
// 			neighbour := Neighbour{
// 				n: wl,
// 			}
// 			neighbours = append(neighbours, neighbour)
// 		}

// 	}
// 	fmt.Println("NEIGHBOURS ARE: ", neighbours)
// 	return neighbours
// }

// // this function find the current year given in those neighbours
// func findDateInNeighbours(year string, nbrs []Neighbour) string {
// 	var date string

// 	for k, v := range nbrs {
// 		if v.n == year {
// 			//date = nbrs[k-2].n + " " + nbrs[k-1].n + " " + nbrs[k].n
// 			date = nbrs[k-2].n + nbrs[k-1].n + nbrs[k].n

// 		}
// 	}

// 	fmt.Println("Date is: ", date)

// 	return date
// }
