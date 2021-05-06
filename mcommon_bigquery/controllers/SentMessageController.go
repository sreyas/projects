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

func GetSentMessages() {
	tableID := "SentMessages"
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
		SentMessageSchema := bigquery.Schema{
			{Name: "MessageID", Type: bigquery.StringFieldType, Required: false},
			{Name: "Type", Type: bigquery.StringFieldType, Required: false},
			{Name: "Status", Type: bigquery.StringFieldType, Required: false},
			{Name: "PhoneNumber", Type: bigquery.StringFieldType, Required: false},
			{Name: "Profile", Type: bigquery.StringFieldType, Required: false},
			{Name: "Body", Type: bigquery.StringFieldType, Required: false},
			{Name: "SentAt", Type: bigquery.StringFieldType, Required: false},
			{Name: "MessageTemplateID", Type: bigquery.StringFieldType, Required: false},
			{Name: "Mms", Type: bigquery.StringFieldType, Required: false},
			{Name: "Multipart", Type: bigquery.StringFieldType, Required: false},
			{Name: "CampaignID", Type: bigquery.StringFieldType, Required: false},
			{Name: "CampaignActive", Type: bigquery.StringFieldType, Required: false},
			{Name: "CampaignName", Type: bigquery.StringFieldType, Required: false},
			{Name: "PreviousID", Type: bigquery.StringFieldType, Required: false},
		}
		metaData := &bigquery.TableMetadata{
			Schema: SentMessageSchema,
		}
		tableRef := bigqueryclient.Dataset(DatasetID).Table(tableID)
		if err := tableRef.Create(ctx, metaData); err != nil {
			log.Println("error noting", err)
		}
	}
	method := "GET"
	campaignUrl := BaseUrl + "campaigns"
	campaignclient := &http.Client{}
	campaignreq, err := http.NewRequest(method, campaignUrl, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	campaignreq.SetBasicAuth(Username, Password)
	campaignres, err := campaignclient.Do(campaignreq)
	if err != nil {
		fmt.Println(err)
		return
	}

	campaignbody, err := ioutil.ReadAll(campaignres.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	campaignres.Body.Close()
	Campaignresponse := CampaignResponse{}
	_ = xml.Unmarshal([]byte(campaignbody), &Campaignresponse)

	for _, campaign := range Campaignresponse.Campaigns.Campaign {
		log.Println(campaign.ID)

		pageNo := 1
		pageCount := 1
		for pageNo <= pageCount {
			TotalSentMessages := []SentMessage{}
			url := BaseUrl + "sent_messages?limit=1000&campaign_id=" + campaign.ID + "&limit=1000&page=" + strconv.Itoa(pageNo) + "&start_time=2021-02-01"

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

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			res.Body.Close()
			messageResponse := SentMessageTotalResponse{}
			_ = xml.Unmarshal([]byte(body), &messageResponse)
			pageCount, _ = strconv.Atoi(messageResponse.Messages.PageCount)
			pageNo += 1
			for _, message := range messageResponse.Messages.Message {
				var singleMessage SentMessage
				singleMessage.MessageID = message.ID
				singleMessage.Type = message.Type
				singleMessage.Status = message.Status
				singleMessage.PhoneNumber = message.PhoneNumber
				singleMessage.Profile = message.Profile
				singleMessage.Body = message.Body
				singleMessage.SentAt = message.SentAt
				singleMessage.MessageTemplateID = message.MessageTemplateID
				singleMessage.Mms = message.Mms
				singleMessage.Multipart = message.Multipart
				singleMessage.CampaignID = message.Campaign.ID
				singleMessage.CampaignActive = message.Campaign.Active
				singleMessage.CampaignName = message.Campaign.Name
				singleMessage.PreviousID = message.PreviousID.ID
				TotalSentMessages = append(TotalSentMessages, singleMessage)

			}
			log.Println("length of sent messages", len(TotalSentMessages))
			inserter := bigqueryclient.Dataset(DatasetID).Table(tableID).Inserter()
			if err := inserter.Put(ctx, TotalSentMessages); err != nil {
				log.Println("insertion error")

			}

		}

	}

}

type SentMessage struct {
	MessageID         string
	Type              string
	Status            string
	PhoneNumber       string
	Profile           string
	Body              string
	SentAt            string
	MessageTemplateID string
	Mms               string
	Multipart         string
	CampaignID        string
	CampaignActive    string
	CampaignName      string
	PreviousID        string
}
type SentMessageResponse struct {
	Text              string `xml:",chardata"`
	ID                string `xml:"id,attr"`
	Type              string `xml:"type,attr"`
	Status            string `xml:"status,attr"`
	PhoneNumber       string `xml:"phone_number"`
	Profile           string `xml:"profile"`
	Body              string `xml:"body"`
	SentAt            string `xml:"sent_at"`
	MessageTemplateID string `xml:"message_template_id"`
	Mms               string `xml:"mms"`
	Multipart         string `xml:"multipart"`
	Campaign          struct {
		Text   string `xml:",chardata"`
		ID     string `xml:"id,attr"`
		Active string `xml:"active,attr"`
		Name   string `xml:"name"`
	} `xml:"campaign"`
	PreviousID struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id,attr"`
	} `xml:"previous_id"`
}
type SentMessageTotalResponse struct {
	XMLName  xml.Name `xml:"response"`
	Text     string   `xml:",chardata"`
	Success  string   `xml:"success,attr"`
	Messages struct {
		Text      string                `xml:",chardata"`
		Page      string                `xml:"page,attr"`
		Limit     string                `xml:"limit,attr"`
		PageCount string                `xml:"page_count,attr"`
		Message   []SentMessageResponse `xml:"message"`
	} `xml:"messages"`
}
