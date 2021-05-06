package listcampaign

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetCampaigns() []Campaign {

	url := "https://secure.mcommons.com/api/campaigns"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Authorization", "Basic Y2lyaWwuc3JlZWRoYXJAc3JleWFzLmNvbTpjaXJpbC5zcmVlZGhhckBzcmV5YXMuY29t")
	req.Header.Add("Cookie", "_mcommons_session_id=fefc59659ed445f1e0eb392413b7a100")

	res, err := client.Do(req)
	defer res.Body.Close()
	campaignResponse := CampaignResponse{}
	body, err := ioutil.ReadAll(res.Body)
	err = xml.Unmarshal(body, &campaignResponse)

	return campaignResponse.Campaigns.Campaign
}

type CampaignResponse struct {
	XMLName   xml.Name `xml:"response"`
	Text      string   `xml:",chardata"`
	Success   string   `xml:"success,attr"`
	Campaigns struct {
		Text     string     `xml:",chardata"`
		Campaign []Campaign `xml:"campaign"`
	} `xml:"campaigns"`
}
type Campaign struct {
	Text        string `xml:",chardata"`
	ID          string `xml:"id,attr"`
	Active      string `xml:"active,attr"`
	Name        string `xml:"name"`
	Description string `xml:"description"`
	Tags        struct {
		Text string `xml:",chardata"`
		Tag  string `xml:"tag"`
	} `xml:"tags"`
}
