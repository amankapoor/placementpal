package scraper

import "time"

import "math/rand"

func genRandInt(len int) int {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return r1.Intn(len)

}

type UserAgent struct {
	ua string
}

var UserAgents = []UserAgent{
	{
		ua: "Mozilla/5.0 (Linux; Android 6.0.1; A0001 Build/MHC19Q) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.83 Mobile Safari/537.36",
	},
	{
		ua: "Mozilla/5.0 (Android 4.4.4; Mobile; rv:56.0) Gecko/56.0 Firefox/56.0",
	},
	{
		ua: "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1",
	},
	{
		ua: "Mozilla/5.0 (iPhone; CPU iPhone OS 10_2 like Mac OS X) AppleWebKit/602.3.12 (KHTML, like Gecko) Version/10.0 Mobile/14C92 Safari/602.1",
	},
	{
		ua: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	},
	{
		ua: "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) CriOS/56.0.2924.75 Mobile/14E5239e Safari/602.1",
	},
}

func RotateUserAgents() string {
	uas := UserAgents
	// fmt.Println(len(uas))
	n := genRandInt(len(uas))
	return uas[n].ua
}
