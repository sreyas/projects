package htmllink

import (
	"fmt"
	"log"

	"net/http"
	"strconv"
	"strings"
)

func ParseHtml(campaignId, broadCastId, sessionId string) BroadCastDetail {

	url := "https://secure.mcommons.com/campaigns/" + campaignId + "/broadcasts/" + broadCastId + "/report"

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Cookie", sessionId)

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if err != nil {
		panic(err)
	}
	// resp, _ := http.Get(url)
	ress := &http.Response{}
	// ress := res.Body
	*ress = *res
	defer res.Body.Close()

	links, err := Parse(res.Body)

	log.Println(res.Cookies())
	if err != nil {
		panic(err)
	}

	broadCast := parseMcommons(links)

	return broadCast
}

type BroadCastDetail struct {
	Receiptients float64
	OPT_OUT      float64
	OPT_OUT_Per  float64
	Reply        float64
	Reply_Per    float64
	WebClicks    float64
	WebClick_Per float64
	WebUrl       string
	Responses    float64
	Response_Per float64
	MessageType  string
}

func parseMcommons(links []Link) BroadCastDetail {
	//Parse Sent Messages
	var broadCast BroadCastDetail
	var err error
	for linkIndex, link := range links {
		if strings.Contains(link.Text, "recipients") {
			ReceiptientsMessages := strings.Split(link.Text, " ")

			broadCast.Receiptients, err = strconv.ParseFloat(strings.Replace(ReceiptientsMessages[0], ",", "", -1), 64)
			if err != nil {
				log.Println(err)
			}

		} else if strings.Contains(link.Text, "recipient") {
			ReceiptientsMessages := strings.Split(link.Text, " ")

			broadCast.Receiptients, err = strconv.ParseFloat(strings.Replace(ReceiptientsMessages[0], ",", "", -1), 64)
			if err != nil {
				log.Println(err)
			}
		}

		if strings.Contains(link.Href, "opt_out") {

			broadCast.OPT_OUT, err = strconv.ParseFloat(strings.Replace(link.Text, ",", "", -1), 64)
			if err != nil {
				log.Println(err)
			}
			if broadCast.OPT_OUT != 0 && broadCast.Receiptients != 0 {
				perc := (broadCast.OPT_OUT / broadCast.Receiptients) * 100
				perc, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", perc), 64)
				broadCast.OPT_OUT_Per = perc
			} else {
				broadCast.OPT_OUT_Per = 0
			}
		}

		if strings.Contains(link.Href, "reply") {

			broadCast.Reply, err = strconv.ParseFloat(strings.Replace(link.Text, ",", "", -1), 64)

			if err != nil {
				log.Println(err)
			}
			if broadCast.Reply != 0 && broadCast.Receiptients != 0 {
				perc := (broadCast.Reply / broadCast.Receiptients) * 100
				perc, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", perc), 64)
				broadCast.Reply_Per = perc
			} else {
				broadCast.Reply_Per = 0
			}

		}
		if strings.Contains(link.Text, "Web Clicks") {
			WebClicks := strings.Split(link.Text, " ")
			broadCast.WebClicks, err = strconv.ParseFloat(strings.Replace(WebClicks[0], ",", "", -1), 64)
			if err != nil {
				log.Println("web clicks err", err)
			}
			log.Println("web clicks err", err)
			broadCast.WebUrl = links[linkIndex+1].Text

			if broadCast.WebClicks != 0 && broadCast.Receiptients != 0 {
				var perc float64
				perc = (broadCast.WebClicks / broadCast.Receiptients) * 100
				log.Println(perc)
				log.Println(broadCast.WebClicks, broadCast.Receiptients)
				perc, err = strconv.ParseFloat(fmt.Sprintf("%.2f", perc), 64)
				log.Println(err)
				broadCast.WebClick_Per = perc
			} else {
				broadCast.Response_Per = 0
			}
		}
		if strings.Contains(link.Text, "responses") {
			ResponsesMessages := strings.Split(link.Text, " ")

			broadCast.Responses, err = strconv.ParseFloat(strings.Replace(ResponsesMessages[0], ",", "", -1), 64)
			if err != nil {
				log.Println(err)
			}
			if broadCast.Responses != 0 && broadCast.Receiptients != 0 {
				perc := (broadCast.Responses / broadCast.Receiptients) * 100
				perc, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", perc), 64)
				broadCast.Response_Per = perc
			} else {
				broadCast.Response_Per = 0
			}

		}
		if strings.Contains(link.Text, "SMS MT") {
			broadCast.MessageType = "SMS"
		} else if strings.Contains(link.Text, "MMS MT") {
			broadCast.MessageType = "MMS"
		}

	}

	return broadCast
}
