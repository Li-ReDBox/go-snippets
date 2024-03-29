package bq

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type Saver interface {
	SaveView(content, name string)
	SaveSchema(content []byte, name string)
}

func CreatDataset(ds *bigquery.Dataset, ctx context.Context, client *bigquery.Client, location string) error {
	meta := &bigquery.DatasetMetadata{
		Location: location, // See https://cloud.google.com/bigquery/docs/locations
	}
	if err := ds.Create(ctx, meta); err != nil {
		return err
	}
	fmt.Println(ds.DatasetID, "has been created successfully.")
	return nil
}

// GetOrCreateDataset gets information of a dataset. If it does not exist, create it first.
func GetOrCreateDataset(datasetID string, ctx context.Context, client *bigquery.Client, location string) (*bigquery.Dataset, error) {
	// this dataset not necessarily exist
	ds := client.Dataset(datasetID)
	meta, err := ds.Metadata(ctx)

	if err == nil {
		fmt.Printf("Dataset ID: %s\n", datasetID)
		fmt.Printf("Description: %s\n", meta.Description)
		fmt.Println("Labels:")
		for k, v := range meta.Labels {
			fmt.Printf("\t%s: %s", k, v)
		}
	} else {
		if !isNotExist(err) {
			return nil, err
		}
		err = CreatDataset(ds, ctx, client, location)
		if err != nil {
			return nil, fmt.Errorf("dataset %s does not exist, and failed to create it: %s", datasetID, err)
		}
		// CreatDataset does not set Description or Lables, so no details need to printed
		fmt.Printf("Created a new Dataset of ID: %s\n", datasetID)
	}

	return ds, nil
}

func NewDatamart(projectID, location, name string) Datamart {
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}

	ds, err := GetOrCreateDataset(name, ctx, client, location)

	if err != nil {
		client.Close()
		log.Fatalln("Cannot get or create dataset", name)
	}

	return Datamart{
		ctx:       ctx,
		client:    client,
		datasetID: name,
		dataset:   ds,
	}
}

func PackageDataset(ctx context.Context, client *bigquery.Client, ds *bigquery.Dataset) Datamart {
	return Datamart{
		ctx:       ctx,
		datasetID: ds.DatasetID,
		client:    client,
		dataset:   ds,
	}
}

type Datamart struct {
	ctx       context.Context
	datasetID string
	client    *bigquery.Client
	dataset   *bigquery.Dataset
}

func (d Datamart) Close() {
	d.client.Close()
}

func (d Datamart) GetTables() []string {
	// this also has code to print labels
	// this dataset has to exist
	tables := []string{}
	meta, err := d.dataset.Metadata(d.ctx)
	if err != nil {
		fmt.Println("Cannot download dataset", d.datasetID, "has error", err)
		return tables
	}

	fmt.Printf("Dataset ID: %s\n", d.datasetID)
	fmt.Printf("Description: %s\n", meta.Description)
	fmt.Println("Labels:")
	for k, v := range meta.Labels {
		fmt.Printf("\t%s: %s", k, v)
	}
	fmt.Println("Tables:")
	it := d.client.Dataset(d.datasetID).Tables(d.ctx)

	cnt := 0
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		cnt++
		fmt.Printf("\t%s\n", t.TableID)
		tables = append(tables, t.TableID)
	}
	if cnt == 0 {
		fmt.Println("\tThis dataset does not contain any tables.")
	}
	return tables
}

func (d Datamart) GetView(viewID string) (string, error) {
	view := d.dataset.Table(viewID)
	meta, err := view.Metadata(d.ctx)
	if err != nil {
		return "", err
	}
	fmt.Printf("View %s, query: %s\n", view.FullyQualifiedName(), meta.ViewQuery)
	return meta.ViewQuery, nil
}

func (d Datamart) GetTable(id string) (bigquery.Schema, error) {
	table := d.dataset.Table(id)
	meta, err := table.Metadata(d.ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("View %s\n", table.FullyQualifiedName())
	return meta.Schema, nil
}

func printSchema(schema bigquery.Schema) {
	for _, s := range schema {
		fmt.Printf("%v\n", s)
		fmt.Printf("%s: %s, %s, %t %t\n\n", s.Name, s.Description, s.Type, s.Required, s.Repeated)
		if s.Type == bigquery.RecordFieldType {
			printSchema(s.Schema)
		}
	}

}

// Download retrieve schemas and send them to a Saver
func (d Datamart) Download(tables []string, saver Saver) {
	var wg sync.WaitGroup

	for _, table := range tables {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()

			bt := d.dataset.Table(id)
			meta, err := bt.Metadata(d.ctx)
			if err != nil {
				fmt.Printf("Cannot get meta of %s of dataset %s, err: %s", d.datasetID, id, err)
				return
			}

			var def string
			if meta.Type == bigquery.RegularTable {
				printSchema(meta.Schema)
				js, err := getSchemaJSON(meta.Schema)
				if err != nil {
					fmt.Printf("Cannot get JSON of the schema of %s of dataset %s, err: %s", d.datasetID, id, err)
					return
				}
				saver.SaveSchema(js, id)
			} else if meta.Type == bigquery.ViewTable {
				def = meta.ViewQuery
				saver.SaveView(def, id)
			} else {
				fmt.Println(bt, "is not supported for downloading", meta.Type)
			}

		}(table)
	}
	wg.Wait()
}

// CreateView creates a view of the query
func (d Datamart) CreateView(query, viewID string) {
	meta := &bigquery.TableMetadata{
		ViewQuery: query,
	}
	if err := d.dataset.Table(viewID).Create(d.ctx, meta); err != nil {
		log.Fatalln("Could not create", viewID, err)
	}
	fmt.Println("View", viewID, "has been created successfully.")
}

func isNotExist(e error) bool {
	es := e.Error()
	return strings.Contains(es, "404")
}
