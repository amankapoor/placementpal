package pdfextractor

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/sajari/docconv"
)

type year struct {
	year string
}

var years = []year{
	{
		year: "2018",
	},
	{
		year: "2019",
	},
	{
		year: "2017",
	},
}

// This will return PDF data struct containing Batch, Title, and Degrees
func DataFromPDF(fileName string) ([]string, []string, error) {
	var res *docconv.Response
	var err error
	//Get the pdf file
	res, err = docconv.ConvertPath("../../cmd/apid/views/temp/" + fileName)
	if err != nil {
		return nil, nil, errors.Wrap(err, "<<Unable to get result from filename>>")
	}

	// Save file body
	s := res.Body
	// Save body also as a slice
	arr := strings.Fields(s)
	//fmt.Println(arr)

	// first we find batches 2017, 2018 and 2019
	var b []string
	for _, v := range years {
		batch := FindBatch(arr, v.year)
		//fmt.Println("BATCH: ", batch)
		b = append(b, batch)
	}
	batches := RemoveWhiteSpacesFromStringSlice(b)

	// WE ARE NOW INSTEAD SHOWING TITLE TAKEN FROM SITE
	// then we find title
	// title := FindTitle(arr)
	// fmt.Println("TITLE: ", title)

	// find degrees
	degrees := FindDegrees(nil, arr)
	//fmt.Println("FIND DEGREE RETURNED: ", degrees)

	return batches, degrees, nil
}
