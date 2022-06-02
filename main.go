package main

import (
	"context"
	"io/ioutil"
	"math/rand"
	"time"

	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {
	bucketName := "savetypictures"
	url := "https://miro.medium.com/max/404/1*65iXGLup5igJDZXA1oOFXw.png" //replace with URL to actual car/id picture
	fileName := RandStringBytes(20)                                       //any unique name will do. I just use random text.
	storeGCS(url, bucketName, fileName)
	text1 := "somedatafield"    //replace with real data
	text2 := "anotherdatafield" //replace with real data
	appendrow(text1, text2, fileName)
}

func appendrow(text1 string, text2 string, fileName string) {
	data, err := ioutil.ReadFile("text-recognition-320818-44031b293769.json")
	checkError(err)
	conf, err := google.JWTConfigFromJSON(data, sheets.SpreadsheetsScope)
	checkError(err)

	client := conf.Client(context.TODO())
	srv, err := sheets.New(client)
	checkError(err)

	spreadsheetID := "1coX_shoKfObvPM_e2rKoNixmtB2cyDxw6lWlj8tWq54"
	readRange := "Database"
	var vr sheets.ValueRange

	var savetydata []interface{}
	savetydata = append(savetydata, text1)
	savetydata = append(savetydata, text2)
	savetydata = append(savetydata, fileName)

	vr.Values = append(vr.Values, savetydata)

	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, readRange, &vr).ValueInputOption("RAW").Do()
	checkError(err)

}

func storeGCS(url, bucketName, fileName string) error {
	// Create GCS connection
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("text-recognition-320818-44031b293769.json"))
	// Connect to bucket
	bucket := client.Bucket(bucketName)
	checkError(err)
	// Get the url response
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP response error: %v", err)
	}
	// response.Body.Close()
	if response.StatusCode == http.StatusOK {
		// Setup the GCS object with the filename to write to
		obj := bucket.Object(fileName)

		// w implements io.Writer.
		w := obj.NewWriter(ctx)

		// Copy file into GCS
		if _, err := io.Copy(w, response.Body); err != nil {
			return fmt.Errorf("Failed to copy to bucket: %v", err)
		}

		// Close, just like writing a file. File appears in GCS after
		if err := w.Close(); err != nil {
			return fmt.Errorf("Failed to close: %v", err)
		}
	}
	return nil
}
