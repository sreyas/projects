package controllers

import (
	"encoding/xml"
)

type ClickBigquery struct {
	ClickID     string
	CreatedAt   string
	URL         string
	ClickedURL  string
	RemoteAddr  string
	HTTPReferer string
	UserAgent   string
	Url_ID      string
	Profile_ID  string
}
type Click struct {
	Text        string `xml:",chardata"`
	ClickID     string `xml:"id,attr"`
	CreatedAt   string `xml:"created_at"`
	URL         string `xml:"url"`
	ClickedURL  string `xml:"clicked_url"`
	RemoteAddr  string `xml:"remote_addr"`
	HTTPReferer string `xml:"http_referer"`
	UserAgent   string `xml:"user_agent"`
	Url_ID      string
	Profile     struct {
		ID string `xml:"id,attr"`
	} `xml:"profile"`
}

type ClickResponse struct {
	XMLName xml.Name `xml:"response"`
	Text    string   `xml:",chardata"`
	Success string   `xml:"success,attr"`
	Clicks  struct {
		Text   string  `xml:",chardata"`
		Num    string  `xml:"num,attr"`
		Unique string  `xml:"unique,attr"`
		Page   string  `xml:"page,attr"`
		Click  []Click `xml:"click"`
	} `xml:"clicks"`
}
