package main

import (
	"context"
	"os"
	"log"
	"fmt"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var applicationCredentialsFile string = os.Getenv("FICHIS_GOOGLE_APPLICATION_CREDENTIALS_FILE_PATH")
var projectID string = os.Getenv("FICHIS_GOOGLE_PROJECT_ID")
var domainName string = os.Getenv("FICHIS_DOMAIN_NAME")

func initFirestoreClient(ctx context.Context) *firestore.Client {
        // Sets your Google Cloud Platform project ID.
        client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(applicationCredentialsFile))
        if err != nil {
                log.Fatalf("Failed to create client: %v", err)
        }
        return client
}
var domainPrefix string = fmt.Sprintf("https://%s/", domainName)
var ctx = context.Background()
var client = initFirestoreClient(ctx)

func addLink(key string, value string) (link string, err error) {
	link = domainPrefix + key
	_, err = client.Collection("urls").Doc(key).Create(ctx, map[string]interface{}{
		"link": value,
		})
	return

}

func deleteLink(key string) (err error) {
	_, err = client.Collection("urls").Doc(key).Delete(ctx)
	return
}

func getLink(key string) (link string, err error) {
	snapshot, err := client.Collection("urls").Doc(key).Get(ctx)
	data := snapshot.Data()
	link = fmt.Sprint(data["link"])
	return
}
