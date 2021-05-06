package controllers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

func GetKeyWords() {
	tableID := "Keywords"
	ctx := context.Background()

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
		KeywordSchema := bigquery.Schema{
			{Name: "KeywordID", Type: bigquery.StringFieldType, Required: false},
			{Name: "Active", Type: bigquery.StringFieldType, Required: false},
			{Name: "Name", Type: bigquery.StringFieldType, Required: false},
			{Name: "OptInPathID", Type: bigquery.StringFieldType, Required: false},
			{Name: "EndedAt", Type: bigquery.StringFieldType, Required: false},
			{Name: "CreatedAt", Type: bigquery.StringFieldType, Required: false},
			{Name: "UpdatedAt", Type: bigquery.StringFieldType, Required: false},
		}
		metaData := &bigquery.TableMetadata{
			Schema: KeywordSchema,
		}
		tableRef := bigqueryclient.Dataset(DatasetID).Table(tableID)
		if err := tableRef.Create(ctx, metaData); err != nil {
			log.Println("error noting", err)
		}
	}
	url := BaseUrl + "keywords"
	method := "GET"

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
	ExistingKeyWords := GetExistingKeyWords(ctx, bigqueryclient)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	res.Body.Close()
	keywordResponse := KeyWordTotalResponse{}
	_ = xml.Unmarshal([]byte(body), &keywordResponse)
	TotalKeyWords := []Keyword{}
	for _, keyword := range keywordResponse.Keywords.Keyword {
		var singleKeyword Keyword
		singleKeyword.KeywordID = keyword.ID
		singleKeyword.Active = keyword.Active
		singleKeyword.Name = keyword.Name
		singleKeyword.OptInPathID = keyword.OptInPathID
		singleKeyword.EndedAt = keyword.EndedAt
		singleKeyword.CreatedAt = keyword.CreatedAt
		singleKeyword.UpdatedAt = keyword.UpdatedAt
		if _, ok := ExistingKeyWords[singleKeyword.KeywordID]; !ok {
			TotalKeyWords = append(TotalKeyWords, singleKeyword)
		} else {
			log.Println("keyword found")
			updateQuery := GenerateUpdateQuery(tableID, singleKeyword, "KeywordID", singleKeyword.KeywordID)
			log.Println(updateQuery)
			query := bigqueryclient.Query(updateQuery)
			_, err := query.Run(ctx)
			log.Println(err)

		}

	}
	// inserter := bigqueryclient.Dataset(DatasetID).Table(tableID).Inserter()
	// if err := inserter.Put(ctx, TotalKeyWords); err != nil {
	// 	log.Println("insertion error")
	// 	log.Println(err)
	// }

}

type Keyword struct {
	KeywordID   string
	Active      string
	Name        string
	OptInPathID string
	EndedAt     string
	CreatedAt   string
	UpdatedAt   string
}
type KeywordResponse struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id,attr"`
	Active      string `xml:"active,attr"`
	Name        string `xml:"name"`
	OptInPathID string `xml:"opt_in_path_id"`
	EndedAt     string `xml:"ended_at"`
	CreatedAt   string `xml:"created_at"`
	UpdatedAt   string `xml:"updated_at"`
}
type KeyWordTotalResponse struct {
	XMLName  xml.Name `xml:"response"`
	Text     string   `xml:",chardata"`
	Success  string   `xml:"success,attr"`
	Keywords struct {
		Text    string            `xml:",chardata"`
		Keyword []KeywordResponse `xml:"keyword"`
	} `xml:"keywords"`
}

func GetExistingKeyWords(ctx context.Context, bigqueryclient *bigquery.Client) map[string]Keyword {
	tableID := "Keywords"
	allKeyword := make(map[string]Keyword)
	query := bigqueryclient.Query(
		`SELECT * FROM .` + DatasetID + `.` + tableID + ``)
	bigqryRes, err := query.Read(ctx)
	log.Println(err)
	for {
		var row Keyword
		err := bigqryRes.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		allKeyword[row.KeywordID] = row
		// fmt.Fprintf(w, "url: %s views: %s\n", row.URL, row.ClickID)
	}
	return allKeyword
}
