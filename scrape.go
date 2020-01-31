package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
)

func main() {
	doc, err := goquery.NewDocument("https://tflare.com/testscrapeconnpass/")
	if err != nil {
		fmt.Print("url scarapping failed")
	}
	doc.Find("div.user_info > a.image_link").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		narrow(url)
	})
}

func narrow(url string) {

	// 管理者 https://connpass.com/user/tflare/open/
	open, _ := regexp.Compile(`^https://connpass.com/user/(.*?)/open/$`)
	// 発表者 https://connpass.com/user/tflare/presentation/
	presentation, _ := regexp.Compile(`^https://connpass.com/user/(.*?)/presentation/$`)
	// 参加者 https://connpass.com/user/tflare/
	attendance, _ := regexp.Compile(`^https://connpass.com/user/(.*?)/$`)

	regDataOpen := open.FindStringSubmatch(url)
	if regDataOpen != nil {
		return;
	}

	regDataPresentation := presentation.FindStringSubmatch(url)
	if regDataPresentation != nil {
		fmt.Println(regDataPresentation[1])
		return;
	}

	regDataAttendance := attendance.FindStringSubmatch(url)
	if regDataAttendance != nil {
		fmt.Println(regDataAttendance[1])
		return;
	}

}
