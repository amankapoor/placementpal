package pdfextractor

import (
	"regexp"
)

type Degree struct {
	Name string
}

func FindDegrees(dgs []Degree, s []string) []string {
	dgs = Degrees
	var availableDegrees []string

	for _, v := range dgs {
		//fmt.Printf("Currently executing k: %d | v: %v \n", k, v)
		deg := findDegree(s, v.Name)
		availableDegrees = append(availableDegrees, deg)
		//fmt.Println("Executed and available degrees equlas", availableDegrees)
	}

	d := RemoveWhiteSpacesFromStringSlice(availableDegrees)
	return d
}

func RemoveWhiteSpacesFromStringSlice(slc []string) []string {
	var r []string
	for _, v := range slc {
		if v != "" {
			r = append(r, v)
		}
	}
	return r
}

func findDegree(s []string, degree string) string {
	var r *regexp.Regexp
	switch degree {
	case "B.Tech":
		r = regexp.MustCompile(`^(\W?)(?i)(b(\.)?(tech)?)(\W?)$`)
	case "BBA":
		r = regexp.MustCompile(`^(\W?)(?i)(b(\.)?b(\.)?a(\.)?)(\W?)$`)
	case "BCA":
		r = regexp.MustCompile(`^(\W?)(?i)(b(\.)?c(\.)?a(\.)?)(\W?)$`)
	case "B.Sc":
		r = regexp.MustCompile(`^(\W?)(?i)(b(\.)?s(\.)?c(\.)?)(\W?)$`)
	// commenting out BE degree because it also matches be english word
	// case "B.E.":
	// 	r = regexp.MustCompile(`^(\W?)(?i)(b(\.)?e(\.)?)(\W?)$`)
	case "B.Com":
		r = regexp.MustCompile(`^(\W?)(?i)(b(\.)?(com)?)(\W?)$`)
	case "BA":
		r = regexp.MustCompile(`^(\W?)(?i)(b(\.)?a(\.)?)(\W?)$`)
	case "B.Arch":
		r = regexp.MustCompile(`^(\W?)(?i)(b(\.)?(arch)?)(\W?)$`)
	case "B.Pharma":
		r = regexp.MustCompile(`^(\W?)(?i)(b(\.)?(pharma)?)(\W?)$`)
	case "BFA":
		r = regexp.MustCompile(`^(\W?)(?i)(b(\.)?f(\.)?a(\.)?)(\W?)$`)
	case "BPT":
		r = regexp.MustCompile(`^(\W?)(?i)(b(\.)?p(\.)?t(\.)?)(\W?)$`)
	case "M.Tech":
		//r = regexp.MustCompile(`^([\W]?(M|m){1}[\.]*[T|t]{1}[e|E]{1}[c|C]{1}[h|H]{1}[\W]?)+$`)
		r = regexp.MustCompile(`^(\W?)(?i)(m(\.)?(tech)?)(\W?)$`)
	case "MBA":
		//r = regexp.MustCompile(`^([\W]?(M|m){1}[\.]*[B|b]{1}[\.]*[A|a]{1}[\W]?)+$`)
		r = regexp.MustCompile(`^(\W?)(?i)(m(\.)?b(\.)?a(\.)?)(\W?)$`)
	case "MCA":
		//r = regexp.MustCompile(`^([\W]?(M|m){1}[\.]*[C|c]{1}[\.]*[A|a]{1}[\W]?)+$`)
		r = regexp.MustCompile(`^(\W?)(?i)(m(\.)?c(\.)?a(\.)?)(\W?)$`)
	case "M.Sc":
		//r = regexp.MustCompile(`^([\W]?(M|m){1}[\.]*[S|s]{1}[\.]*[c|C]{1}[\W]?)+$`)
		r = regexp.MustCompile(`^(\W?)(?i)(m(\.)?s(\.)?c(\.)?)(\W?)$`)
	case "M.Pharma":
		r = regexp.MustCompile(`^(\W?)(?i)(m(\.)?(pharma)?)(\W?)$`)
	case "MPT":
		r = regexp.MustCompile(`^(\W?)(?i)(m(\.)?p(\.)?t(\.)?)(\W?)$`)
	case "MPH":
		r = regexp.MustCompile(`^(\W?)(?i)(m(\.)?p(\.)?h(\.)?)(\W?)$`)
	case "D.Pharmacy":
		r = regexp.MustCompile(`^(\W?)(?i)(d(\.)?(pharmacy)?)(\W?)$`)
	default:
		return ""
	}

	var tf bool
	//var st []string
	var deg string
	//fmt.Println("Value of r is: ", r)

	locs := findDegreeWordLocations(s)
	nbrs := findNeighboursOfDegreeWord(locs, s)
	for _, v := range nbrs {
		tf = r.MatchString(v.n)
		//fmt.Printf("k: %d | v: (%v) | st: %v\n", k, v.n, tf)
		if tf == true {
			deg = degree
			//fmt.Println("Degree returned,", deg)
			return deg
		}
	}
	return ""
}

// this function tries to find the location of word 'batch' in the pdf
func findDegreeWordLocations(body []string) []Location {
	var count int
	locations := make([]Location, 0, 10)
	r := regexp.MustCompile(`^(\W?)(?i)(degrees?)(\W?)$`)
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

	// fmt.Println("Count is", count)
	// fmt.Println("Location are, ", locations)
	return locations
}

type NeighbourDegree struct {
	n string
}

// this function find the neighbours of the word batch so that we can find the number
func findNeighboursOfDegreeWord(loc []Location, s []string) []NeighbourDegree {
	neighbours := make([]NeighbourDegree, 0, 82)
	for i := 0; i < len(loc); i++ {
		pos := loc[i].loc

		for j := pos - 20; j <= pos+20; j++ {
			wl := findWordAtPosition(j, s)
			neighbour := NeighbourDegree{
				n: wl,
			}
			neighbours = append(neighbours, neighbour)
		}
	}
	// fmt.Println("NEIGHBOURS ARE: ", neighbours)
	return neighbours
}

var Degrees = []Degree{
	{
		Name: "B.Tech",
	},
	{
		Name: "BBA",
	},
	{
		Name: "BCA",
	},
	{
		Name: "B.Sc",
	},
	{
		Name: "B.Com",
	},
	{
		Name: "BA",
	},
	{
		Name: "B.Arch",
	},
	{
		Name: "B.Pharma",
	},
	{
		Name: "BFA",
	},
	{
		Name: "BPT",
	},
	{
		Name: "M.Tech",
	},
	{
		Name: "MBA",
	},
	{
		Name: "MCA",
	},
	{
		Name: "M.Sc",
	},
	{
		Name: "M.Pharma",
	},
	{
		Name: "MPT",
	},
	{
		Name: "MPH",
	},
	{
		Name: "D.Pharmacy",
	},
}

// FOR FUTURE USE
// var Branches = []Branch{
// 	{
// 		Name: "CSE",
// 	},
// 	{
// 		Name: "IT",
// 	},
// 	{
// 		Name: "ECE",
// 	},
// 	{
// 		Name: "ET",
// 	},
// 	{
// 		Name: "EI",
// 	},
// 	{
// 		Name: "EEE",
// 	},
// 	{
// 		Name: "MAE",
// 	},
// 	{
// 		Name: "Civil",
// 	},
// 	{
// 		Name: "Aerospace",
// 	},
// 	{
// 		Name: "Biotech",
// 	},
// 	{
// 		Name: "Bioinfo",
// 	},
// 	{
// 		Name: "ICE",
// 	},
// 	{
// 		Name: "Foodtech",
// 	},
// 	{
// 		Name: "Nanotech",
// 	},
// 	{
// 		Name: "Pharma",
// 	},
// 	{
// 		Name: "Avionics",
// 	},
// 	{
// 		Name: "Computer Application",
// 	},
// 	{
// 		Name: "Automobile",
// 	},
// 	{
// 		Name: "Embedded Systems",
// 	},
// 	{
// 		Name: "Mechatronics",
// 	},
// 	{
// 		Name: "VLSI",
// 	},
// 	{
// 		Name: "Power Systems",
// 	},
// 	{
// 		Name: "Solar & Alternative Energy",
// 	},
// 	{
// 		Name: "NTM",
// 	},
// 	{
// 		Name: "Physics",
// 	},
// 	{
// 		Name: "Chemistry",
// 	},
// 	{
// 		Name: "Mathematics",
// 	},
// 	{
// 		Name: "Computer Application (MCA)",
// 	},
// }

// FOR FUTURE USE
// var Institutes = []Institute{
// 	{
// 		Name: "ASET",
// 	},
// 	{
// 		Name: "AIAE",
// 	},
// 	{
// 		Name: "AIB",
// 	},
// 	{
// 		Name: "AIFT",
// 	},
// 	{
// 		Name: "AINT",
// 	},
// 	{
// 		Name: "AISST",
// 	},
// 	{
// 		Name: "AITTM",
// 	},
// 	{
// 		Name: "AIIT",
// 	},
// 	{
// 		Name: "ASCS",
// 	},
// 	{
// 		Name: "AIRAE",
// 	},
// 	{
// 		Name: "ASE",
// 	},
// }

// // FOR FUTURE USE
// var InstituteLocations = []InstituteLocation{
// 	{
// 		Name: "Noida",
// 	},
// 	{
// 		Name: "Delhi",
// 	},
// 	{
// 		Name: "Lucknow",
// 	},
// 	{
// 		Name: "Jaipur",
// 	},
// 	{
// 		Name: "Gurgaon",
// 	},
// 	{
// 		Name: "Gwalior",
// 	},
// 	{
// 		Name: "Greater Noida",
// 	},
// }
