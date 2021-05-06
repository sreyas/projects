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

type MessageResponse struct {
	Text              string `xml:",chardata"`
	ID                string `xml:"id,attr"`
	Type              string `xml:"type,attr"`
	Approved          string `xml:"approved,attr"`
	PhoneNumber       string `xml:"phone_number"`
	CarrierName       string `xml:"carrier_name"`
	Profile           string `xml:"profile"`
	Body              string `xml:"body"`
	MessageTemplateID string `xml:"message_template_id"`
	Mms               string `xml:"mms"`
	Multipart         string `xml:"multipart"`
	Keyword           string `xml:"keyword"`
	ReceivedAt        string `xml:"received_at"`
	PreviousID        struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id,attr"`
	} `xml:"previous_id"`
	Campaign struct {
		Text   string `xml:",chardata"`
		ID     string `xml:"id,attr"`
		Active string `xml:"active,attr"`
		Name   string `xml:"name"`
	} `xml:"campaign"`
	NextID struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id,attr"`
	} `xml:"next_id"`
}
type MessageTotalResponse struct {
	XMLName  xml.Name `xml:"response"`
	Text     string   `xml:",chardata"`
	Success  string   `xml:"success,attr"`
	Messages struct {
		Text      string            `xml:",chardata"`
		Page      string            `xml:"page,attr"`
		Limit     string            `xml:"limit,attr"`
		PageCount string            `xml:"page_count,attr"`
		Message   []MessageResponse `xml:"message"`
	} `xml:"messages"`
}
type Message struct {
	MessageID         string
	Type              string
	Approved          string
	PhoneNumber       string
	CarrierName       string
	Profile           string
	Body              string
	MessageTemplateID string
	Mms               string
	Multipart         string
	Keyword           string
	ReceivedAt        string
	PreviousID        string
	NextID            string
	CampaignID        string
	CampaignName      string
}

func GetMessages() {
	tableID := "Messages"
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
		MessageSchema := bigquery.Schema{
			{Name: "MessageID", Type: bigquery.StringFieldType, Required: false},
			{Name: "Type", Type: bigquery.StringFieldType, Required: false},
			{Name: "Approved", Type: bigquery.StringFieldType, Required: false},
			{Name: "PhoneNumber", Type: bigquery.StringFieldType, Required: false},
			{Name: "CarrierName", Type: bigquery.StringFieldType, Required: false},
			{Name: "Profile", Type: bigquery.StringFieldType, Required: false},
			{Name: "Body", Type: bigquery.StringFieldType, Required: false},
			{Name: "MessageTemplateID", Type: bigquery.StringFieldType, Required: false},
			{Name: "Mms", Type: bigquery.StringFieldType, Required: false},
			{Name: "Multipart", Type: bigquery.StringFieldType, Required: false},
			{Name: "Keyword", Type: bigquery.StringFieldType, Required: false},
			{Name: "ReceivedAt", Type: bigquery.StringFieldType, Required: false},
			{Name: "PreviousID", Type: bigquery.StringFieldType, Required: false},
			{Name: "NextID", Type: bigquery.StringFieldType, Required: false},
			{Name: "CampaignID", Type: bigquery.StringFieldType, Required: false},
			{Name: "CampaignName", Type: bigquery.StringFieldType, Required: false},
		}
		metaData := &bigquery.TableMetadata{
			Schema: MessageSchema,
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
			url := BaseUrl + "messages?limit=1000&campaign_id=" + campaign.ID + "&limit=1000&page=" + strconv.Itoa(pageNo)

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
			messageResponse := MessageTotalResponse{}
			_ = xml.Unmarshal([]byte(body), &messageResponse)
			pageCount, _ = strconv.Atoi(messageResponse.Messages.PageCount)
			pageNo += 1
			TotalMessages := []Message{}
			for _, message := range messageResponse.Messages.Message {
				var singleMessage Message
				singleMessage.MessageID = message.ID
				singleMessage.Type = message.Type
				singleMessage.Approved = message.Approved
				singleMessage.PhoneNumber = message.PhoneNumber
				singleMessage.CarrierName = message.CarrierName
				singleMessage.Profile = message.Profile
				singleMessage.Body = message.Body
				singleMessage.MessageTemplateID = message.MessageTemplateID
				singleMessage.Mms = message.Mms
				singleMessage.Multipart = message.Multipart
				singleMessage.Keyword = message.Keyword
				singleMessage.ReceivedAt = message.ReceivedAt
				singleMessage.PreviousID = message.PreviousID.ID
				singleMessage.NextID = message.NextID.ID
				singleMessage.CampaignID = message.Campaign.ID
				singleMessage.CampaignName = message.Campaign.Name

				TotalMessages = append(TotalMessages, singleMessage)

			}
			log.Println(len(TotalMessages))
			inserter := bigqueryclient.Dataset(DatasetID).Table(tableID).Inserter()
			if err := inserter.Put(ctx, TotalMessages); err != nil {
				log.Println("insertion error")
				log.Println(err)
			}
		}

	}

}

type CampaignResponse struct {
	XMLName   xml.Name `xml:"response"`
	Text      string   `xml:",chardata"`
	Success   string   `xml:"success,attr"`
	Campaigns struct {
		Text     string `xml:",chardata"`
		Campaign []struct {
			Text        string `xml:",chardata"`
			ID          string `xml:"id,attr"`
			Active      string `xml:"active,attr"`
			Name        string `xml:"name"`
			Description string `xml:"description"`
			Tags        struct {
				Text string `xml:",chardata"`
				Tag  string `xml:"tag"`
			} `xml:"tags"`
		} `xml:"campaign"`
	} `xml:"campaigns"`
}
