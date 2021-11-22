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
	flag.Parse()

	if *folder == "" {
		log.Fatalln("Missing folder for saving queries. Make sure it has been created.")
	}

	fmt.Println("About to save", *datasetID, " of project set by the project set by env variable to folder", *folder)
	d := bq.NewDatamart(*datasetID)
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
