package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func GetClicksDetails() {

	if ProjectID == "" {
		fmt.Println("GOOGLE_CLOUD_PROJECT environment variable must be set.")
		os.Exit(1)
	}

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, ProjectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	rows, err := query(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	if err := printResults(os.Stdout, rows); err != nil {
		log.Fatal(err)
	}
}

// query returns a row iterator suitable for reading query results.
func query(ctx context.Context, client *bigquery.Client) (*bigquery.RowIterator, error) {
	tableId := "Clicks"
	query := client.Query(
		`SELECT * FROM .` + DatasetID + `.` + tableId + ``)
	return query.Read(ctx)
}

// printResults prints results from a query to the Stack Overflow public dataset.
func printResults(w io.Writer, iter *bigquery.RowIterator) error {
	CLicks := []ClickBigquery{}
	for {
		var row ClickBigquery
		err := iter.Next(&row)
		if err == iterator.Done {
			log.Println(len(CLicks))
			return nil
		}
		if err != nil {
			return fmt.Errorf("error iterating through results: %v", err)
		}
		CLicks = append(CLicks, row)
		// fmt.Fprintf(w, "url: %s views: %s\n", row.URL, row.ClickID)
	}

	return nil
}
