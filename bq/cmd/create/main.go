package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"

	"funmech/bq"
)

func main() {
	var datasetID = flag.String("d", "li_test_go", "name of datatset to be created")
	var sourceDir = flag.String("f", "", "name of folder holds sql files")
	flag.Parse()

	projectID := os.Getenv("TARGET_PROJECTID")
	if projectID == "" {
		log.Fatalln("GOOGLE_CLOUD_PROJECT environment variable must be set.")
	}

	fmt.Println("About to create Dataset", *datasetID, "in project", projectID)

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Cannot create bigquery.NewClient: %v", err)
	}

	ds, err := bq.GetOrCreateDataset(*datasetID, ctx, client)
	if err != nil {
		fmt.Printf("Unhandled error when getting or creating dataset %s, %+v", *datasetID, err)
	}

	datamart := bq.PackageDataset(ctx, client, ds)
	defer datamart.Close()

	createAll(*sourceDir, "*.sql", datamart.CreateView)
}
