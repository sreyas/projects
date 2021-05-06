package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

const (
	ProjectID = ""
	DatasetID = ""
	Username  = ""
	Password  = ""
	BaseUrl   = "https://secure.mcommons.com/api/"
)

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

type ConfigFile struct {
	Clicks struct {
		LastDate string `json: "lastdate"`
	} `json: "clicks"`
	Profile struct {
		LastDate string `json:"lastdate"`
	} `json:"profile"`
}

func ReadConfFile() ConfigFile {
	file, _ := ioutil.ReadFile("conf.json")

	confData := ConfigFile{}

	_ = json.Unmarshal([]byte(file), &confData)
	return confData
}
func GenerateConfigTables() {
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
	query := client.Query(
		`SELECT * FROM .` + DatasetID + `.` + tableID + ``)
	res, err := query.Read(ctx)
	Configs := []JobDetailsTable{}
	for {
		var row JobDetailsTable
		err := res.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("error iterating through results: %v", err)
		}
		Configs = append(Configs, row)
		// fmt.Fprintf(w, "url: %s views: %s\n", row.URL, row.ClickID)
	}
	log.Println(Configs)
}
func GetTableQueryParams(tableId string) (string, string) {
	jobdetails := GetJobDetails(tableId)
	layout := "2006-01-02"
	var startDate, endDate string
	if jobdetails.LastRun != "" {

		startDateTime, _ := time.Parse(layout, jobdetails.LastRun)
		startDate = startDateTime.AddDate(0, 0, 1).Format("2006-01-02")

	}
	endDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	return startDate, endDate
}
func GenerateUpdateQuery(destTableId string, Values interface{}, checkingname string, CheckingValue string) string {
	updateqry := ""
	updateqry += fmt.Sprintf("update .%s.%s set ", DatasetID, destTableId)
	v := reflect.ValueOf(Values)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		limiter := ","
		if i == v.NumField()-1 {
			limiter = ""
		}
		updateqry += fmt.Sprintf(` %s="%s" %s`, typeOfS.Field(i).Name, v.Field(i).Interface(), limiter)
	}
	updateqry += fmt.Sprintf(` where %s="%s" `, checkingname, CheckingValue)
	return updateqry
}
