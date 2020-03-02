package main

import (
	"context"

	"fmt"
	"os"

	"cloud.google.com/go/datastore"
)

func main() {
	const sourceKind = "GeoAddressNormalised"
	const batchSize = 500

	ctx := context.Background()

	// Create a datastore client. In a typical application, you would create
	// a single client which is reused for every datastore operation.
	client, err := datastore.NewClient(context.Background(), datastore.DetectProjectID)
	if err != nil {
		fmt.Printf("datastore.NewClient: %v. Missing GOOGLE_APPLICATION_CREDENTIALS?\n", err)
		os.Exit(1)
	}

	fmt.Printf("About to delete all entities of %s\n", sourceKind)
	query := datastore.NewQuery(sourceKind).KeysOnly()

	keys, err := client.GetAll(ctx, query, nil)
	if err != nil {
		fmt.Println("Cannot get all keys", err)
		return
	}

	count := len(keys)
	if count > batchSize {
		start, end := 0, 0
		for end < count-1 {
			end = start + batchSize - 1
			if end > count {
				end = count - 1
				fmt.Println("Set end to count", end)
			}
			if err := client.DeleteMulti(ctx, keys[start:end]); err != nil {
				fmt.Println("Fail to delete", err)
				break
			} else {
				fmt.Printf("Deleted records from %d to %d\n", start, end)
				start = end + 1
			}
		}
	} else {
		if err := client.DeleteMulti(ctx, keys); err != nil {
			fmt.Println("Fail to delete", err)
		} else {
			fmt.Printf("Deleted %d records\n", len(keys))
		}
	}
}
