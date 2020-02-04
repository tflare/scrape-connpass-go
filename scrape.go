package main

import (
	"context"
	"log"
	"regexp"

	"github.com/PuerkitoBio/goquery"

	"google.golang.org/api/option"

	"cloud.google.com/go/firestore"
)

var projectID = "attendance-functions"

func main() {

	ctx := context.Background()
	sa := option.WithCredentialsFile("/home/tflare/attendance-functions-b1c2438d620c.json")
	client, err := firestore.NewClient(ctx, projectID, sa)
	if err != nil {
		log.Fatalf("Failtd to create client: %v", err)
	}
	defer client.Close()

	scrape(ctx, client)
}

func scrape(ctx context.Context, client *firestore.Client) {
	doc, err := goquery.NewDocument("https://tflare.com/testscrapeconnpass/")
	if err != nil {
		log.Fatalf("url scarapping failed")
	}

	doc.Find("div.user_info > a.image_link").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		narrow(ctx, client, url)
	})
}

func narrow(ctx context.Context, client *firestore.Client, url string) {

	// 管理者 https://connpass.com/user/tflare/open/
	open, _ := regexp.Compile(`^https://connpass.com/user/(.*?)/open/$`)
	// 発表者 https://connpass.com/user/tflare/presentation/
	presentation, _ := regexp.Compile(`^https://connpass.com/user/(.*?)/presentation/$`)
	// 参加者 https://connpass.com/user/tflare/
	attendance, _ := regexp.Compile(`^https://connpass.com/user/(.*?)/$`)

	regDataOpen := open.FindStringSubmatch(url)
	if regDataOpen != nil {
		return
	}

	regDataPresentation := presentation.FindStringSubmatch(url)
	if regDataPresentation != nil {
		writeDB(ctx, client, regDataPresentation[1], true)
		return
	}

	regDataAttendance := attendance.FindStringSubmatch(url)
	if regDataAttendance != nil {
		writeDB(ctx, client, regDataAttendance[1], false)
		return
	}

}

func writeDB(ctx context.Context, client *firestore.Client, userID string, presenter bool) {
	_, _, err := client.Collection("attendance").Add(ctx, map[string]interface{}{
		"eventID":    151286,
		"userID":     userID,
		"attendance": false, //出席フラグ今の段階ではfalseで登録
		"presenter":  presenter,
		"createdAt":  firestore.ServerTimestamp,
		"updatedAt":  firestore.ServerTimestamp,
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
}
