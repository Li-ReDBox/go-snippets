package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"

	"funmech/bq"
)

func main() {
	var datasetID = flag.String("dataset", "li_test_go", "name of datatset to be created")
	var sourceDir = flag.String("folder", "", "name of folder holds sql files")
	var projectID = flag.String("project", "", "Project ID where BQ is")
	var location = flag.String("location", "", "Location of BQ")
	flag.Parse()

	bq.Required(projectID, "TARGET_PROJECTID", "Project ID where BQ is")
	bq.Required(location, "BQ_LOCATION", "Location of BQ")

	fmt.Println("About to create Dataset", *datasetID, "in project", projectID)

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, *projectID)
	if err != nil {
		log.Fatalf("Cannot create bigquery.NewClient: %v", err)
	}

	ds, err := bq.GetOrCreateDataset(*datasetID, ctx, client, *location)
	if err != nil {
		fmt.Printf("Unhandled error when getting or creating dataset %s, %+v", *datasetID, err)
	}

	datamart := bq.PackageDataset(ctx, client, ds)
	defer datamart.Close()

	createAll(*sourceDir, "*.sql", datamart.CreateView)
}
