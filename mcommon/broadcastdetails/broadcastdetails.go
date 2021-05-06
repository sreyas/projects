package broadcastdetails

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetBroadcast(broadCastId string) BroadCast {

	url := "https://secure.mcommons.com/api/broadcast?broadcast_id=" + broadCastId
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Basic "+basicAuth("", ""))

	// req.Header.Add("Authorization", "Basic Y2lyaWwuc3JlZWRoYXJAc3JleWFzLmNvbTowcGNsWW1McFlOdDE=")
	// req.Header.Add("Cookie", "_mcommons_session_id=fefc59659ed445f1e0eb392413b7a100")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	broadcastResponse := BroadCastResponse{}
	err = xml.Unmarshal(body, &broadcastResponse)
	return broadcastResponse.Broadcast
}
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type BroadCastResponse struct {
	XMLName   xml.Name  `xml:"response"`
	Text      string    `xml:",chardata"`
	Success   string    `xml:"success,attr"`
	Broadcast BroadCast `xml:"broadcast"`
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
	RepliesCount             string `xml:"replies_count"`
	OptOutsCount             string `xml:"opt_outs_count"`
	EstimatedRecipientsCount string `xml:"estimated_recipients_count"`
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
	Tags           string `xml:"tags"`
}
