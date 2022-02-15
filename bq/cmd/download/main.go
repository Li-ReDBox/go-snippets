package main

import (
	"flag"
	"fmt"
	"log"

	"funmech/bq"
)

func main() {
	var datasetID = flag.String("d", "li_test_go", "name of datatset")
	var folder = flag.String("f", "", "name of folder to save queries")
	var projectID = flag.String("project", "", "Project ID where BQ is")
	var location = flag.String("location", "", "Location of BQ")
	flag.Parse()

	if *folder == "" {
		log.Fatalln("Missing folder for saving queries. Make sure it has been created.")
	}

	bq.Required(projectID, "TARGET_PROJECTID", "Project ID where BQ is")
	bq.Required(location, "BQ_LOCATION", "Location of BQ")

	fmt.Println("About to save", *datasetID, " of project set by the project set by env variable to folder", *folder)
	// Assume a dataset only contains RegularTable or ViewTable
	d := bq.NewDatamart(*projectID, *location, *datasetID)
	defer d.Close()

	// a bit of debugging print
	tables := d.GetTables()
	for k, v := range tables {
		fmt.Println(k+1, v)
	}

	// save table defintions to local files system
	lf := LocalStorage{
		Folder: *folder,
	}
	d.Download(tables, lf)
}
