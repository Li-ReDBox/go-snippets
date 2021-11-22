package bq

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type Saver interface {
	Save(path, name string)
}

func CreatDataset(ds *bigquery.Dataset, ctx context.Context, client *bigquery.Client) error {
	meta := &bigquery.DatasetMetadata{
		Location: os.Getenv("BQ_LOCATION"), // See https://cloud.google.com/bigquery/docs/locations
	}
	if err := ds.Create(ctx, meta); err != nil {
		return err
	}
	fmt.Println(ds.DatasetID, "has been created successfully.")
	return nil
}

// GetOrCreateDataset gets information of a dataset. If it does not exist, create it first.
func GetOrCreateDataset(datasetID string, ctx context.Context, client *bigquery.Client) (*bigquery.Dataset, error) {
	// this dataset not necessarily exist
	ds := client.Dataset(datasetID)
	meta, err := ds.Metadata(ctx)
	if err != nil {
		if !isNotExist(err) {
			return nil, err
		}
		err = CreatDataset(ds, ctx, client)
		if err != nil {
			return nil, fmt.Errorf("Cannot create dataset %s", err)
		}
	}

	meta, err = ds.Metadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("Cannot get meta of dataset %s, %s", datasetID, err)
	}

	fmt.Printf("Dataset ID: %s\n", datasetID)
	fmt.Printf("Description: %s\n", meta.Description)
	fmt.Println("Labels:")
	for k, v := range meta.Labels {
		fmt.Printf("\t%s: %s", k, v)
	}
	return ds, nil
}

func NewDatamart(name string) Datamart {
	projectID := os.Getenv("TARGET_PROJECTID")
	if projectID == "" {
		fmt.Println("GOOGLE_CLOUD_PROJECT environment variable must be set.")
		os.Exit(1)
	}

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}

	ds, err := GetOrCreateDataset(name, ctx, client)

	if err != nil {
		client.Close()
		fmt.Println("Cannot get or create dataset", name)
		os.Exit(1)
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

// Download retrieve schemas and send them to a Saver
func (d Datamart) Download(tables []string, saver Saver) {
	var wg sync.WaitGroup

	for _, t := range tables {
		wg.Add(1)
		go func(view string) {
			defer wg.Done()

			def, err := d.GetView(view)
			if err != nil {
				fmt.Printf("Cannot view %s of dataset %s, err: %s", d.datasetID, view, err)
			}
			saver.Save(def, view)
		}(t)
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
