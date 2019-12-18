package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"cloud.google.com/go/datastore"
)

type LinkMap struct {
	TargetURL string `datastore:",noindex"`
	Handle    string
	GAPage    string
	ItemID    int
}

func main() {
	ctx := context.Background()

	// pretent Datastore entity is set from environment variable
	os.Setenv("DS_ENTITY", "Book")
	var dsClient *datastore.Client

	dsClient, err := datastore.NewClient(context.Background(), datastore.DetectProjectID)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get all entries
	var entities []LinkMap
	keys, err := dsClient.GetAll(ctx, datastore.NewQuery(os.Getenv("DS_ENTITY")), &entities)
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, key := range keys {
		fmt.Println(key)
		fmt.Println(entities[i])
	}

	export, err := json.Marshal(entities)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Json object")
		fmt.Println(string(export))
		// & in urls are encoded to \u0026. Not fixed. Plain text version of it is in backup_noencode.json
		if err := ioutil.WriteFile("backup.json", export, 0666); err != nil {
			fmt.Println(err)
		}
	}

	// // Delete all retrieved entries
	// if err := dsClient.DeleteMulti(ctx, keys); err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Should have deleted %d entries.", len(keys))
}
