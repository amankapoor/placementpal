package pdfextractor

import (
	"strings"
)

// findTitle takes in the full pdf body as string
// and return the title by stopping at the first occurence
// of the word "Batch" as this word was common across all the pdfs
// Then it returns everything it finds before the word "Batch"
func FindTitle(s []string) string {

	title, exists := throughFirstBatchOccurence(s)
	//fmt.Printf("Title is %s and exists is %v", title, exists)
	if exists == true {
		t := removeAmityPlacementCode(title)
		return t
	}

	// fallback
	t := firstNWordsAsTitle(7, s)
	return t
}

func throughFirstBatchOccurence(s []string) (string, bool) {
	var stoppedAt int
	var t []string
	var title string
	for k, v := range s {
		if v == "Batch" || v == "batch" {
			stoppedAt = k
			//fmt.Println("k is %d & v is %v", k, v)
			t = s[:stoppedAt+1]
			title = strings.Join(t, " ")
			//fmt.Println(title)
			return title, true
		}
	}

	return "", false
}

//this is fallback
func firstNWordsAsTitle(n int, s []string) string {
	str := s[:7]
	title := strings.Join(str, " ")
	return title
}

// removeAmityPlacementCode() takes in the title parsed from findTitle()
// and breaks it into array. Then it analyses if the first element is
// Amity's placement code like TSC17180058. If true, it return a new title
// without placement code, else returns the original string
func removeAmityPlacementCode(title string) string {
	s := strings.Fields(title)
	contains := strings.Contains(s[0], "1718")
	//fmt.Println("PALCEMENT CODE: \n", contains)
	if contains {
		return strings.Join(s[1:], " ")
	}
	return title
}
