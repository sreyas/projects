package listcampaignbroadcasts

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetBroadCasts() []Broadcast {

	url := "https://secure.mcommons.com/api/broadcasts?limit=1000"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Basic "+basicAuth("", ""))
	res, err := client.Do(req)
	defer res.Body.Close()
	broadCastResponse := BroadCastResponse{}
	body, err := ioutil.ReadAll(res.Body)
	err = xml.Unmarshal(body, &broadCastResponse)
	// log.Println(string(body))
	var cookievalue string
	cookies := res.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "_mcommons_session_id" {
			cookievalue = cookie.Value
		}
	}
	if cookievalue != "" {

	}
	return broadCastResponse.Broadcasts.Broadcast

}

type BroadCastResponse struct {
	XMLName    xml.Name `xml:"response"`
	Text       string   `xml:",chardata"`
	Success    string   `xml:"success,attr"`
	Broadcasts struct {
		Text      string      `xml:",chardata"`
		Page      string      `xml:"page,attr"`
		Limit     string      `xml:"limit,attr"`
		PageCount string      `xml:"page_count,attr"`
		Broadcast []Broadcast `xml:"broadcast"`
	} `xml:"broadcasts"`
}

type Broadcast struct {
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
		Group []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
			ID   string `xml:"id,attr"`
			Type string `xml:"type,attr"`
		} `xml:"group"`
	} `xml:"included_groups"`
	ExcludedGroups struct {
		Text  string `xml:",chardata"`
		Group []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
			ID   string `xml:"id,attr"`
			Type string `xml:"type,attr"`
		} `xml:"group"`
	} `xml:"excluded_groups"`
	Tags struct {
		Text string `xml:",chardata"`
		Tag  string `xml:"tag"`
	} `xml:"tags"`
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
