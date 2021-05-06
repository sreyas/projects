package controllers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func getExistingTinyUrls(ctx context.Context, bigqueryclient *bigquery.Client) map[string]TinyURL {
	tableID := "TinyUrls"
	allTinyurs := make(map[string]TinyURL)
	query := bigqueryclient.Query(
		`SELECT * FROM .` + DatasetID + `.` + tableID + ``)
	bigqryRes, err := query.Read(ctx)
	log.Println(err)
	for {
		var row TinyURL
		err := bigqryRes.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		allTinyurs[row.TinyURLID] = row
		// fmt.Fprintf(w, "url: %s views: %s\n", row.URL, row.ClickID)
	}
	return allTinyurs
}
func GetTinyURLS() {
	tableID := "TinyUrls"
	clickTableId := "Clicks"
	ctx := context.Background()

	bigqueryclient, err := bigquery.NewClient(ctx, ProjectID)

	defer bigqueryclient.Close()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	tableFound := false
	clickTableFound := false
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
		if t.TableID == clickTableId {
			clickTableFound = true
		}

	}
	if !tableFound {

		TinyUrlSchema := bigquery.Schema{
			{Name: "Name", Type: bigquery.StringFieldType, Required: false},
			{Name: "Mode", Type: bigquery.StringFieldType, Required: false},
			{Name: "Url", Type: bigquery.StringFieldType, Required: false},
			{Name: "Host", Type: bigquery.StringFieldType, Required: false},
			{Name: "Description", Type: bigquery.StringFieldType, Required: false},
			{Name: "Key", Type: bigquery.StringFieldType, Required: false},
			{Name: "Created", Type: bigquery.StringFieldType, Required: false},
			{Name: "TinyURLID", Type: bigquery.StringFieldType, Required: false},
		}
		metaData := &bigquery.TableMetadata{
			Schema: TinyUrlSchema,
		}
		tableRef := bigqueryclient.Dataset(DatasetID).Table(tableID)
		if err := tableRef.Create(ctx, metaData); err != nil {
			log.Println("error noting", err)
		}
	}
	if !clickTableFound {
		clickTableSchema := bigquery.Schema{

			{Name: "ClickID", Type: bigquery.StringFieldType, Required: false},
			{Name: "CreatedAt", Type: bigquery.StringFieldType, Required: false},
			{Name: "URL", Type: bigquery.StringFieldType, Required: false},
			{Name: "ClickedURL", Type: bigquery.StringFieldType, Required: false},
			{Name: "RemoteAddr", Type: bigquery.StringFieldType, Required: false},
			{Name: "HTTPReferer", Type: bigquery.StringFieldType, Required: false},
			{Name: "UserAgent", Type: bigquery.StringFieldType, Required: false},
			{Name: "Url_ID", Type: bigquery.StringFieldType, Required: false},
			{Name: "Profile_ID", Type: bigquery.StringFieldType, Required: false},
		}
		metaData := &bigquery.TableMetadata{
			Schema: clickTableSchema,
		}
		tableRef := bigqueryclient.Dataset(DatasetID).Table(clickTableId)
		if err := tableRef.Create(ctx, metaData); err != nil {
			log.Println("error noting", err)
		}
	}

	url := "https://secure.mcommons.com/api/tinyurls"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println("94", err)
		return
	}
	req.SetBasicAuth(Username, Password)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("100", err)
		return
	}
	currentTinyUrls := getExistingTinyUrls(ctx, bigqueryclient)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	res.Body.Close()
	TinyURLs := TinyURLTotalResponse{}
	_ = xml.Unmarshal([]byte(body), &TinyURLs)
	TotalTinyUrls := []TinyURL{}
	var wg sync.WaitGroup
	wg.Add(1)
	CounterLimit := 10
	TinyURLs.Tinyurls.Tinyurl = TinyURLs.Tinyurls.Tinyurl[1:2]
	log.Println("aishdsauhd", len(TinyURLs.Tinyurls.Tinyurl))
	limit := len(TinyURLs.Tinyurls.Tinyurl)

	startDate, endDate := GetTableQueryParams(clickTableId)

	for tinyUrlIndex, tinyUrl := range TinyURLs.Tinyurls.Tinyurl {
		_ = tinyUrlIndex
		var SingleTinyUrl TinyURL
		SingleTinyUrl.Name = tinyUrl.Name
		SingleTinyUrl.Host = tinyUrl.Host
		SingleTinyUrl.Created = tinyUrl.CreatedAt
		SingleTinyUrl.Description = tinyUrl.Description
		SingleTinyUrl.Key = tinyUrl.Key
		SingleTinyUrl.TinyURLID = tinyUrl.ID
		SingleTinyUrl.Mode = tinyUrl.Mode
		SingleTinyUrl.Url = tinyUrl.URL
		if _, ok := currentTinyUrls[SingleTinyUrl.TinyURLID]; !ok {
			TotalTinyUrls = append(TotalTinyUrls, SingleTinyUrl)
		} else {
			updateQuery := fmt.Sprintf(`update .%s.%s set Name="%s",Host="%s",Created="%s",Description="%s",Key="%s",Mode="%s",Url="%s" where TinyURLID="%s"`, DatasetID, tableID, SingleTinyUrl.Name, SingleTinyUrl.Host, SingleTinyUrl.Created, SingleTinyUrl.Description, SingleTinyUrl.Key, SingleTinyUrl.Mode, SingleTinyUrl.Url, SingleTinyUrl.TinyURLID)
			UpdateTinyUrl(ctx, bigqueryclient, updateQuery)
		}

		go GetClicksController(SingleTinyUrl, &wg, clickTableId, bigqueryclient, &ctx, 1, startDate, endDate)

		if tinyUrlIndex%CounterLimit == 0 && tinyUrlIndex != 0 {
			log.Println("waiting", tinyUrlIndex)

			wg.Wait()
			log.Println("completed", tinyUrlIndex)
			if (limit - tinyUrlIndex) < CounterLimit {
				log.Println("dlkahsdiluashbdjlisag")
				log.Println(limit - tinyUrlIndex)
				wg.Add(limit - tinyUrlIndex - 1)
			} else {
				wg.Add(CounterLimit)
			}
		} else if tinyUrlIndex == 0 {
			log.Println("kausndkasndunuadiod")
			wg.Wait()
			if (limit - tinyUrlIndex) < CounterLimit {
				log.Println("dlkahsdiluashbdjlisag")
				log.Println(limit - tinyUrlIndex)
				wg.Add(limit - tinyUrlIndex - 1)
			} else {
				log.Println("asdbasbu")
				wg.Add(CounterLimit)
			}
		}
		// clicksinserter := bigqueryclient.Dataset(DatasetID).Table(clickTableId).Inserter()
		// if err := clicksinserter.Put(ctx, clicks); err != nil {
		// 	log.Println("insertion error")
		// 	log.Println(err)
		// }

	}
	wg.Wait()
	//Insertion
	UpdateJobDetails(ctx, bigqueryclient, clickTableId, endDate)
	inserter := bigqueryclient.Dataset(DatasetID).Table(tableID).Inserter()
	inserter.SkipInvalidRows = true
	if err := inserter.Put(ctx, TotalTinyUrls); err != nil {
		log.Println("insertion error")
		log.Println(err)
	}

}
func (i *TinyURL) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"Name":        i.Name,
		"Mode":        i.Mode,
		"Url":         i.Url,
		"Host":        i.Host,
		"Description": i.Description,
		"Key":         i.Key,
		"TinyURLID":   i.TinyURLID,
		"Created":     i.Created,
	}, bigquery.NoDedupeID, nil
}
func (i *ClickBigquery) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"ClickID":     i.ClickID,
		"CreatedAt":   i.CreatedAt,
		"URL":         i.URL,
		"ClickedURL":  i.ClickedURL,
		"RemoteAddr":  i.RemoteAddr,
		"HTTPReferer": i.HTTPReferer,
		"UserAgent":   i.UserAgent,
		"Url_ID":      i.Url_ID,
		"Profile_ID":  i.Profile_ID,
	}, bigquery.NoDedupeID, nil
}

type TinyURL struct {
	Name        string
	Mode        string
	Url         string
	Host        string
	Description string
	Key         string
	Created     string
	TinyURLID   string
}
type TinyURLResponse struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id,attr"`
	CreatedAt   string `xml:"created_at"`
	Name        string `xml:"name"`
	Mode        string `xml:"mode"`
	URL         string `xml:"url"`
	Host        string `xml:"host"`
	Description string `xml:"description"`
	Key         string `xml:"key"`
}
type TinyURLTotalResponse struct {
	XMLName  xml.Name `xml:"response"`
	Text     string   `xml:",chardata"`
	Success  string   `xml:"success,attr"`
	Tinyurls struct {
		Text    string            `xml:",chardata"`
		Num     string            `xml:"num,attr"`
		Tinyurl []TinyURLResponse `xml:"tinyurl"`
	} `xml:"tinyurls"`
}

func GetClicksController(TinyURL TinyURL, wg *sync.WaitGroup, clickTableId string, bigqueryclient *bigquery.Client, ctx *context.Context, page int, startDate string, endDate string) {
	defer wg.Done()
	// configDetails := ReadConfFile()
	// lastRunDate := configDetails.Clicks.LastDate
	// layout := "2006-01-02"
	// startDateTime, _ := time.Parse(layout, lastRunDate)
	// startDate := startDateTime.AddDate(0, 0, 1).Format("2006-01-02")
	// toDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	method := "GET"

	TotalClicks := []ClickBigquery{}

	num := 1000

	for num == 1000 {

		clickUrl := "https://secure.mcommons.com/api/clicks?url_id=" + TinyURL.TinyURLID + "&limit=1000&include_profile=true&page=" + strconv.Itoa(page) + "&from=" + startDate + "&to=" + endDate
		log.Println(TinyURL.TinyURLID, " ", page)
		log.Println(clickUrl)
		page += 1

		clickReq, reqerr := http.NewRequest(method, clickUrl, nil)
		if reqerr != nil {
			fmt.Println("71", reqerr)
			// return nil
		}
		client := &http.Client{}
		clickReq.SetBasicAuth(Username, Password)
		clickres, Doerr := client.Do(clickReq)
		if Doerr != nil {
			fmt.Println("78", Doerr)
			// return nil
		}
		clickbody, readAllerr := ioutil.ReadAll(clickres.Body)
		if readAllerr != nil {
			fmt.Println("83", readAllerr.Error())
			log.Println(clickUrl)
			// return nil
			log.Println("breaked")
			GetClicksController(TinyURL, wg, clickTableId, bigqueryclient, ctx, page, startDate, endDate)
			break
		}

		ClickResponse := ClickResponse{}
		xmlerr := xml.Unmarshal([]byte(clickbody), &ClickResponse)
		clickres.Body.Close()
		log.Println("89", xmlerr)
		num, _ = strconv.Atoi(ClickResponse.Clicks.Num)
		for _, v := range ClickResponse.Clicks.Click {

			var singleClick ClickBigquery
			singleClick.ClickID = v.ClickID
			singleClick.CreatedAt = v.CreatedAt
			singleClick.URL = v.URL
			singleClick.ClickedURL = v.ClickedURL
			singleClick.RemoteAddr = v.RemoteAddr
			singleClick.HTTPReferer = v.HTTPReferer
			singleClick.UserAgent = v.UserAgent
			singleClick.Url_ID = TinyURL.TinyURLID
			singleClick.Profile_ID = v.Profile.ID
			TotalClicks = append(TotalClicks, singleClick)
		}

	}
	limit := 10000
	log.Println(TinyURL.TinyURLID, "clicks length", len(TotalClicks))
	for i := 0; i < len(TotalClicks); i += limit {
		clickbatch := TotalClicks[i:min(i+limit, len(TotalClicks))]
		clicksinserter := bigqueryclient.Dataset(DatasetID).Table(clickTableId).Inserter()
		clicksinserter.SkipInvalidRows = true
		if err := clicksinserter.Put(*ctx, clickbatch); err != nil {
			log.Println("insertion error", TinyURL.TinyURLID)
			log.Printf("%+v", clickbatch)
			log.Println(err)
		}
	}

}
func UpdateTinyUrl(ctx context.Context, bigqueryclient *bigquery.Client, updateQuery string) {
	query := bigqueryclient.Query(updateQuery)
	_, _ = query.Read(ctx)

}
