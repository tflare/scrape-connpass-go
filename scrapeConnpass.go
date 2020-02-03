// Package scrapeconnpass provides a set of Cloud Functions.
package scrapeconnpass

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"google.golang.org/api/option"

	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

var projectID = "attendance-functions"

// ToDo Cloud Function で動作するように修正

// scrapeConnpass is an HTTP Cloud Function with a request parameter.
func scrapeConnpass(response http.ResponseWriter, request *http.Request) {
	ctx := context.Background()

	auth(ctx, request)

	sa := option.WithCredentialsFile("/home/tflare/attendance-functions-b1c2438d620c.json")
	client, err := firestore.NewClient(ctx, projectID, sa)
	if err != nil {
		log.Fatalf("Failtd to create client: %v", err)
	}
	defer client.Close()

	scrape(ctx, client)
}

func auth(ctx context.Context, request *http.Request) {

	idToken := getIDToken(request)

	conf := &firebase.Config{
		ServiceAccountID: "",
		ProjectID:        "",
	}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	_, err = authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}
}

func getIDToken(request *http.Request) string {

	// ToDo 動作確認
	reqToken := request.Header.Get("Authorization")
	idToken := strings.Split(reqToken, "Bearer ")
	return idToken[1]
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
