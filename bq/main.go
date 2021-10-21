package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func main() {
	var datasetID = flag.String("d", "li_test", "Dataset name, default=li_test")
	flag.Parse()

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
	defer client.Close()

	ds, err := getOrCreateDataset(*datasetID, ctx, client)
	if err != nil {
		fmt.Printf("Unhandled error when getting or creating dataset %s, %+v", *datasetID, err)
	}
	// This a view of the shakespeare sample dataset, which
	// provides word frequency information.  This view restricts the results to only contain
	// results for works that contain the "king" in the title, e.g. King Lear, King Henry V, etc.
	demoQuery := "SELECT word, word_count, corpus, corpus_date FROM `bigquery-public-data.samples.shakespeare` WHERE corpus LIKE '%king%'"

	createView("simple", demoQuery, ctx, client, ds)

	// rows, err := query(ctx, client)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := printResults(os.Stdout, rows); err != nil {
	// 	log.Fatal(err)
	// }
}

// query returns a row iterator suitable for reading query results.
func query(ctx context.Context, client *bigquery.Client) (*bigquery.RowIterator, error) {

	query := client.Query(
		`SELECT
					CONCAT(
							'https://stackoverflow.com/questions/',
							CAST(id as STRING)) as url,
					view_count
			FROM ` + "`bigquery-public-data.stackoverflow.posts_questions`" + `
			WHERE tags like '%google-bigquery%'
			ORDER BY view_count DESC
			LIMIT 10;`)
	return query.Read(ctx)
}

type StackOverflowRow struct {
	URL       string `bigquery:"url"`
	ViewCount int64  `bigquery:"view_count"`
}

// printResults prints results from a query to the Stack Overflow public dataset.
func printResults(w io.Writer, iter *bigquery.RowIterator) error {
	for {
		var row StackOverflowRow
		err := iter.Next(&row)
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error iterating through results: %v", err)
		}

		fmt.Fprintf(w, "url: %s views: %d\n", row.URL, row.ViewCount)
	}
}

func creatDataset(ds *bigquery.Dataset, ctx context.Context, client *bigquery.Client) {
	meta := &bigquery.DatasetMetadata{
		Location: os.Getenv("BQ_LOCATION"), // See https://cloud.google.com/bigquery/docs/locations
	}
	if err := ds.Create(ctx, meta); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(ds.DatasetID, "has been created successfully.")
}

// printDatasetInfo demonstrates fetching dataset metadata and printing some of it to an io.Writer.
func getOrCreateDataset(datasetID string, ctx context.Context, client *bigquery.Client) (*bigquery.Dataset, error) {
	// this dataset not necessarily exist
	ds := client.Dataset(datasetID)
	meta, err := ds.Metadata(ctx)
	if err != nil {
		if !isNotExist(err) {
			return nil, err
		}
		creatDataset(ds, ctx, client)
		meta, err = ds.Metadata(ctx)
	}

	fmt.Printf("Dataset ID: %s\n", datasetID)
	fmt.Printf("Description: %s\n", meta.Description)
	fmt.Println("Labels:")
	for k, v := range meta.Labels {
		fmt.Printf("\t%s: %s", k, v)
	}
	fmt.Println("Tables:")
	it := client.Dataset(datasetID).Tables(ctx)

	cnt := 0
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		cnt++
		fmt.Printf("\t%s\n", t.TableID)
	}
	if cnt == 0 {
		fmt.Println("\tThis dataset does not contain any tables.")
	}
	return ds, nil
}

// createView creates a view of the query
func createView(query, viewID string, ctx context.Context, client *bigquery.Client, ds *bigquery.Dataset) {
	meta := &bigquery.TableMetadata{
		ViewQuery: query,
	}
	if err := ds.Table(viewID).Create(ctx, meta); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("View", viewID, "has been created successfully.")
}

// getView demonstrates fetching the metadata from a BigQuery logical view and printing it to an io.Writer.
func getView(viewID string, ctx context.Context, client *bigquery.Client, ds *bigquery.Dataset) error {
	// projectID := "my-project-id"
	// datasetID := "mydataset"
	// viewID := "myview"
	view := ds.Table(viewID)
	meta, err := view.Metadata(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("View %s, query: %s\n", view.FullyQualifiedName(), meta.ViewQuery)
	return nil
}

func isNotExist(e error) bool {
	es := e.Error()
	fmt.Println("Error to be checked", es)
	return strings.Contains(es, "404")
}
