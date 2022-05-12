package main

import (
	"context"
	"os"
	"cloud.google.com/go/firestore"
)

func initFirestoreClient(ctx context.Context) *firestore.Client {
        // Sets your Google Cloud Platform project ID.
        projectID := "YOUR_PROJECT_ID"

        client, err := firestore.NewClient(ctx, projectID)
        if err != nil {
                log.Fatalf("Failed to create client: %v", err)
        }
        // Close client when done with
        // defer client.Close()
        return client
}
var domainName string = "https://fich.is/"
var ctx = context.Background()
var client = initFirestoreClient(ctx)
func addLink(key string, value string) (link string, err error) {
	link = domainName + key
	_, _, err := client.Collection("urls").Doc(key).Add(ctx, map[string]interface{}{
		"link": value,
		})
	return

}
func deleteLink(key string) (err error) {
	_, _, err = client.Collection("urls").Doc(key).Delete(ctx)
	return
}

func getLink(key string) (link string, err error) {
}

// _, _, err := client.Collection("users").Add(ctx, map[string]interface{}{
//         "first": "Ada",
//         "last":  "Lovelace",
//         "born":  1815,
// })
