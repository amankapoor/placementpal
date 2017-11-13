package files

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/asciimoo/colly"
	"github.com/pkg/errors"
)

var filespath, temppath string

func init() {
	filespath = filepath.Join("/home/amankapoor/go/src/github.com/amankapoor/placementpal/cmd/apid/files/", "pdf")
	err := os.MkdirAll(filespath, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}

	temppath = filepath.Join("/home/amankapoor/go/src/github.com/amankapoor/placementpal/cmd/apid/views/", "temp")
	err = os.MkdirAll(temppath, os.ModePerm)
	if err != nil {
		log.Println(err)
		return
	}
}

func Downloader(Eid int) ([]string, error) {
	url := "http://amity.edu/placement/Popup.asp?Eid=" + strconv.Itoa(Eid)
	urls := scrapeLinks(url)
	// fmt.Println("URLS are: ", urls)

	var fileName string
	var fileNames []string
	var err error
	for _, v := range urls {
		fileName, err = fileDownloader(v)
		if err != nil {
			return nil, errors.Wrapf(err, "<<unable to download file named: %s>>", fileName)
		}
		fileNames = append(fileNames, fileName)
	}
	return fileNames, nil
}

func fileDownloader(url string) (string, error) {

	u := "http://amity.edu/placement/pdf/" + url

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := netClient.Get(u)
	if err != nil {
		return "", errors.Wrap(err, "<<Unable to initialise client>>")
	}

	if response.StatusCode != 200 {
		return "", errors.Wrap(err, "<<file url did not return 200 status>>")
	}

	if response.ContentLength < 100 {
		return "", errors.Wrap(err, "<<file content length is less than 100 bytes>>")
	}

	defer response.Body.Close()

	//create timestamped filename
	splittedFile := strings.Split(url, ".")
	l := len(splittedFile)
	fileNameWithTimeStamp := splittedFile[l-l] + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + "." + splittedFile[l-1]
	//fmt.Println(fileNameWithTimeStamp)

	// Create output file
	newFile, err := os.Create(strings.Join([]string{temppath, fileNameWithTimeStamp}, "/"))
	if err != nil {
		return "", errors.Wrap(err, "<<unable to create output file>>")
	}
	defer newFile.Close()

	// Write bytes from HTTP response to file.
	// response.Body satisfies the reader interface.
	// newFile satisfies the writer interface.
	// That allows us to use io.Copy which accepts
	// any type that implements reader and writer interface
	_, err = io.Copy(newFile, response.Body)
	if err != nil {
		return "", errors.Wrap(err, "<<Unable to write response body to file>>")
	}
	//log.Printf("Downloaded %d byte file.\n", numBytesWritten)
	return fileNameWithTimeStamp, nil
}

func scrapeLinks(url string) []string {
	var urls []string
	c := colly.NewCollector()
	c.MaxDepth = 1
	c.AllowedDomains = []string{"amity.edu", "localhost:8080"}
	// On every a element which has href attribute call callback
	// c.OnHTML("body table tbody tr td li a[href]", func(e *colly.HTMLElement) {
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

		link := e.Attr("href")
		// Print link
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page

		// Only those links are visited which are matched by  any of the URLFilter regexps
		//c.Visit(e.Request.AbsoluteURL(link))
		a := strings.TrimPrefix(link, "pdf\\")
		urls = append(urls, a)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL.String())
	})

	c.Visit(url)
	//fmt.Println("Placements are: ", placements)
	return urls
}
