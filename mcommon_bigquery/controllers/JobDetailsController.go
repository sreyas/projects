package controllers

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type JobDetailsTable struct {
	TableId string
	LastRun string
}

func GetJobDetails(destTable string) JobDetailsTable {
	tableID := "JobDetails"
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
	log.Println(`SELECT * FROM .` + DatasetID + `.` + tableID + ` where TableId=` + destTable + ``)
	query := client.Query(
		`SELECT * FROM .` + DatasetID + `.` + tableID + ` where TableId="` + destTable + `"`)

	res, err := query.Read(ctx)
	log.Println(err)
	Jobs := []JobDetailsTable{}
	for {
		var row JobDetailsTable
		err := res.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("error iterating through results: %v", err)
		}
		Jobs = append(Jobs, row)
		// fmt.Fprintf(w, "url: %s views: %s\n", row.URL, row.ClickID)
	}
	log.Println(Jobs)
	if len(Jobs) > 0 {
		return Jobs[0]
	} else {
		return JobDetailsTable{}
	}
}

func UpdateJobDetails(ctx context.Context, bigqueryclient *bigquery.Client, destTableId string, endDate string) {
	tableId := "JobDetails"

	updateQuery := fmt.Sprintf(`update .%s.%s set LastRun="%s" where TableId="%s"`, DatasetID, tableId, endDate, destTableId)

	query := bigqueryclient.Query(updateQuery)
	res, err := query.Read(ctx)

	log.Println("jaysgduasydguaysgdujasydgh")
	log.Println(err)
	if err == nil {
		if res.TotalRows == 0 {
			var job JobDetailsTable
			job.TableId = destTableId
			job.LastRun = endDate
			JobInserter := bigqueryclient.Dataset(DatasetID).Table(tableId).Inserter()
			if err := JobInserter.Put(ctx, job); err != nil {
				log.Println("insertion error")
				log.Println(err)

			}
		}
	}
	log.Printf("%+v", res)

}
