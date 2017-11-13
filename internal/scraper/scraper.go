package scraper

import (
	"sort"
	"strconv"
	"strings"

	"github.com/asciimoo/colly"
	"github.com/pkg/errors"
)

// Scrape returns beautified and ordered data after scraping
func FetchLatest() []BeautifiedPlacementData {
	return beautifyData(scrapeAll())
}

// ScrapeAll function scrapes amity.edu/placement,
// upcoming, nontech, internships etc links on placement site
// and returns AllPlacementData of type [][]Placement
func scrapeAll() [][]placementData {
	var allPlacementData [][]placementData

	// // localhost test
	// techMarqueeCSSSelector := "a[href]"
	// techMarqueePlacements := scrapeLinks("http://localhost:5500/dashboard.html", techMarqueeCSSSelector)
	// fmt.Println("TECH MARQUEEE PLACEMENTS: ", techMarqueePlacements)

	techMarqueeCSSSelector := "[style='background:#F6C008; border:#000000; color:#FFFFFF; MARGIN: 3px 3px 3px 3px; font-family:Georgia, serif; font-size:12px ; float:left ; WIDTH: 550px;'] marquee a[href]"
	techMarqueePlacements := scrapeLinks("http://amity.edu/placement/", techMarqueeCSSSelector)
	//fmt.Println("TECH MARQUEEE PLACEMENTS: ", techMarqueePlacements)

	nonTechMarqueeCSSSelector := "[style='background:#01579D; border:#000000; color:#FFFFFF; MARGIN: 3px 3px 3px 3px; font-family:Georgia, serif; font-size:12px ; float:left ; WIDTH: 550px;'] marquee a[href]"
	nonTechMarqueePlacements := scrapeLinks("http://amity.edu/placement/", nonTechMarqueeCSSSelector)
	// fmt.Println("NON TECH MARQUEEE PLACEMENTS: ", nonTechMarqueePlacements)

	// noticeMarqueeCSSSelector := "[style='background:#8d6776; border:#000000; color:#FFFFFF; MARGIN: 3px 3px 3px 3px; font-family:Georgia, serif; font-size:12px ; float:left ; WIDTH: 550px;'] marquee a[href]"
	// noticeMarqueeNotices := scrapeLinks("http://amity.edu/placement/", noticeMarqueeCSSSelector)
	//fmt.Println(noticeMarqueeNotices)

	techRecruitmentsCSSSelector := "#2018 li:nth-child(-n+10) a[href]"
	techRecruitments := scrapeLinks("http://amity.edu/placement/Upcoming_Recu.asp", techRecruitmentsCSSSelector)
	// //fmt.Println(techRecruitments)

	nonTechRecruitmentsCSSSelector := "#2018 li:nth-child(-n+3) a[href]"
	nonTechRecruitments := scrapeLinks("http://amity.edu/placement/Non-Tech_Recu.asp", nonTechRecruitmentsCSSSelector)
	// //fmt.Println(nonTechRecruitments)

	internshipsCSSSelector := "#2018 li a[href]"
	internships := scrapeLinks("http://amity.edu/placement/Internships.asp", internshipsCSSSelector)
	//fmt.Println(internships)

	// allPlacementData = append(allPlacementData, techMarqueePlacements, nonTechMarqueePlacements, techRecruitments, nonTechRecruitments, internships)

	allPlacementData = append(allPlacementData, techMarqueePlacements, nonTechMarqueePlacements, techRecruitments, nonTechRecruitments, internships)
	//fmt.Printf("ALL PLACEMENTS: %v\n", allPlacementData)
	return allPlacementData
}

// scrapeLinks function assists ScrapeAll function by scraping
// the specified url using the css selector. It returns []Placement.
func scrapeLinks(url string, cssSelector string) []placementData {
	var placements []placementData
	c := colly.NewCollector()
	c.UserAgent = RotateUserAgents()
	c.MaxDepth = 1
	c.AllowedDomains = []string{"localhost", "amity.edu"}
	// On every a element which has href attribute call callback
	c.OnHTML(cssSelector, func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		placement := placementData{
			title: e.Text,
			url:   link,
		}
		placements = append(placements, placement)
		// Only those links are visited which are matched by  any of the URLFilter regexps
		//c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(url)
	//fmt.Println("Placements are: ", placements)
	return placements
}

// beautifyData function takes in the output of ScrapeAll function
// and returns a single slice of Placement type instead of
// double slice of Placement that it took as input.
// The output of this function is what we choose to work with.
func beautifyData(p [][]placementData) []BeautifiedPlacementData {
	var bp BeautifiedPlacementData
	var bps []BeautifiedPlacementData
	//var sorted []Placement

	//fmt.Println(len(p))
	for k := range p {
		for _, iv := range p[k] {
			trimmedURL := strings.TrimPrefix(iv.url, "Popup.asp?Eid=")
			urlAsInt, err := strconv.Atoi(trimmedURL)
			if err != nil {
				errors.Wrap(err, "unable to convert atoi")
				// log.Println("Cannot conver atoi at ", trimmedURL)
				// panic("Cannot convert atoi")
			}
			bp = BeautifiedPlacementData{
				Title: iv.title,
				URL:   urlAsInt,
			}
			bps = append(bps, bp)
		}
	}
	sort.SliceStable(bps, func(i, j int) bool {
		sorted := bps[i].URL < bps[j].URL
		return sorted
	})

	return bps
}
