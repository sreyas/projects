package controllers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func GetBroadCastDetails() {

	ctx := context.Background()

	// projectID := ""
	// datasetID := ""
	tableID := "Broadcast"
	bigqueryclient, err := bigquery.NewClient(ctx, ProjectID)
	defer bigqueryclient.Close()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	tableFound := false

	ts := bigqueryclient.Dataset(DatasetID).Tables(ctx)
	for {
		t, err := ts.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {

		}
		if t.TableID == tableID {
			tableFound = true
		}

	}
	if !tableFound {
		BroadCastSchema := bigquery.Schema{
			{Name: "Name", Type: bigquery.StringFieldType, Required: false},
			{Name: "BroadCastId", Type: bigquery.StringFieldType, Required: false},
			{Name: "Body", Type: bigquery.StringFieldType, Required: false},
			{Name: "Campaign_Id", Type: bigquery.StringFieldType, Required: false},
			{Name: "Campaign_Name", Type: bigquery.StringFieldType, Required: false},
			{Name: "DeliveryTime", Type: bigquery.StringFieldType, Required: false},
			{Name: "RepliesCount", Type: bigquery.StringFieldType, Required: false},
			{Name: "OptOutsCount", Type: bigquery.StringFieldType, Required: false},
			{Name: "Include_Subscribers", Type: bigquery.StringFieldType, Required: false},
			{Name: "Throttled", Type: bigquery.StringFieldType, Required: false},
			{Name: "Localtime", Type: bigquery.StringFieldType, Required: false},
			{Name: "Automated", Type: bigquery.StringFieldType, Required: false},
			{Name: "Estimated_Recipients_Count", Type: bigquery.StringFieldType, Required: false},
		}
		metaData := &bigquery.TableMetadata{
			Schema: BroadCastSchema,
		}
		tableRef := bigqueryclient.Dataset(DatasetID).Table(tableID)
		if err := tableRef.Create(ctx, metaData); err != nil {
			log.Println("error noting", err)
		}
	}

	pageNo := 1
	pageCount := 1
	for pageNo <= pageCount {
		url := "https://secure.mcommons.com/api/broadcasts?limit=100&page=" + strconv.Itoa(pageNo)
		method := "GET"
		log.Println(url)
		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
			return
		}
		req.SetBasicAuth(Username, Password)
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		existingBroadCastDetails := GetExistingBroadCastDetails(ctx, bigqueryclient)
		// jsonData, _ := json.Marshal(body)
		BroadCastDetails := BroadCastXml{}
		// json.Unmarshal(jsonData, &BroadCastDetails)
		// log.Println(BroadCastDetails)
		_ = xml.Unmarshal([]byte(body), &BroadCastDetails)
		pageCount, _ = strconv.Atoi(BroadCastDetails.Broadcasts.PageCount)
		pageNo += 1
		TotalBroadCast := []BigqueryBroadCast{}

		for _, v := range BroadCastDetails.Broadcasts.Broadcast {
			var SingleBroadcast BigqueryBroadCast
			SingleBroadcast.Name = v.Name
			SingleBroadcast.BroadCastId = v.ID
			SingleBroadcast.Body = v.Body
			SingleBroadcast.Campaign_Id = v.Campaign.ID
			SingleBroadcast.Campaign_Name = v.Campaign.Name
			SingleBroadcast.DeliveryTime = v.DeliveryTime
			SingleBroadcast.RepliesCount = v.RepliesCount
			SingleBroadcast.OptOutsCount = v.OptOutsCount
			SingleBroadcast.Include_Subscribers = v.IncludeSubscribers
			SingleBroadcast.Throttled = v.Throttled
			SingleBroadcast.Localtime = v.Localtime
			SingleBroadcast.Automated = v.Automated
			SingleBroadcast.Estimated_Recipients_Count = v.EstimatedRecipientsCount
			if _, ok := existingBroadCastDetails[SingleBroadcast.BroadCastId]; !ok {
				TotalBroadCast = append(TotalBroadCast, SingleBroadcast)
			} else {

				updateQuery := GenerateUpdateQuery(tableID, SingleBroadcast, "BroadCastId", SingleBroadcast.BroadCastId)
				query := bigqueryclient.Query(updateQuery)
				_, err := query.Run(ctx)
				log.Println(err)

			}

		}

		// Insertion
		inserter := bigqueryclient.Dataset(DatasetID).Table(tableID).Inserter()
		if err := inserter.Put(ctx, TotalBroadCast); err != nil {
			log.Println("insertion error")
			log.Println(err)
		}
	}

}

type BigqueryBroadCast struct {
	Name                       string
	BroadCastId                string
	Body                       string
	Campaign_Id                string
	Campaign_Name              string
	DeliveryTime               string
	RepliesCount               string
	OptOutsCount               string
	Include_Subscribers        string
	Throttled                  string
	Localtime                  string
	Automated                  string
	Estimated_Recipients_Count string
}

func (i *BigqueryBroadCast) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{

		"Name":                       i.Name,
		"BroadCastId":                i.BroadCastId,
		"Body":                       i.Body,
		"Campaign_Id":                i.Campaign_Id,
		"Campaign_Name":              i.Campaign_Name,
		"DeliveryTime":               i.DeliveryTime,
		"RepliesCount":               i.RepliesCount,
		"OptOutsCount":               i.OptOutsCount,
		"Include_Subscribers":        i.Include_Subscribers,
		"Throttled":                  i.Throttled,
		"Localtime":                  i.Localtime,
		"Automated":                  i.Automated,
		"Estimated_Recipients_Count": i.Estimated_Recipients_Count,
	}, bigquery.NoDedupeID, nil
}

type BroadCast struct {
	Text     string `xml:",chardata"`
	ID       string `xml:"id,attr"`
	Status   string `xml:"status,attr"`
	Name     string `xml:"name"`
	Body     string `xml:"body"`
	Campaign struct {
		Text   string `xml:",chardata"`
		ID     string `xml:"id,attr"`
		Active string `xml:"active,attr"`
		Name   string `xml:"name"`
	} `xml:"campaign"`
	DeliveryTime             string `xml:"delivery_time"`
	IncludeSubscribers       string `xml:"include_subscribers"`
	Throttled                string `xml:"throttled"`
	Localtime                string `xml:"localtime"`
	Automated                string `xml:"automated"`
	EstimatedRecipientsCount string `xml:"estimated_recipients_count"`
	RepliesCount             string `xml:"replies_count"`
	OptOutsCount             string `xml:"opt_outs_count"`
	IncludedGroups           struct {
		Text  string `xml:",chardata"`
		Group struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
			ID   string `xml:"id,attr"`
			Type string `xml:"type,attr"`
		} `xml:"group"`
	} `xml:"included_groups"`
	ExcludedGroups string `xml:"excluded_groups"`
	Tags           struct {
		Text string `xml:",chardata"`
		Tag  string `xml:"tag"`
	} `xml:"tags"`
}

type BroadCastXml struct {
	XMLName    xml.Name `xml:"response"`
	Text       string   `xml:",chardata"`
	Success    string   `xml:"success,attr"`
	Broadcasts struct {
		Text      string      `xml:",chardata"`
		Page      string      `xml:"page,attr"`
		Limit     string      `xml:"limit,attr"`
		PageCount string      `xml:"page_count,attr"`
		Broadcast []BroadCast `xml:"broadcast"`
	} `xml:"broadcasts"`
}

func GetExistingBroadCastDetails(ctx context.Context, bigqueryclient *bigquery.Client) map[string]BigqueryBroadCast {
	tableID := "Broadcast"
	allBroadCast := make(map[string]BigqueryBroadCast)
	query := bigqueryclient.Query(
		`SELECT * FROM .` + DatasetID + `.` + tableID + ``)
	bigqryRes, err := query.Read(ctx)
	log.Println(err)
	for {
		var row BigqueryBroadCast
		err := bigqryRes.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		allBroadCast[row.BroadCastId] = row
		// fmt.Fprintf(w, "url: %s views: %s\n", row.URL, row.ClickID)
	}
	return allBroadCast
}
