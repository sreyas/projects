package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InstagraMResponse struct {
	LoggingPageID         string `json:"logging_page_id"`
	ShowSuggestedProfiles bool   `json:"show_suggested_profiles"`
	ShowFollowDialog      bool   `json:"show_follow_dialog"`
	Graphql               struct {
		User struct {
			Biography              string      `json:"biography"`
			BlockedByViewer        bool        `json:"blocked_by_viewer"`
			RestrictedByViewer     interface{} `json:"restricted_by_viewer"`
			CountryBlock           bool        `json:"country_block"`
			ExternalURL            interface{} `json:"external_url"`
			ExternalURLLinkshimmed interface{} `json:"external_url_linkshimmed"`
			EdgeFollowedBy         struct {
				Count int `json:"count"`
			} `json:"edge_followed_by"`
		} `json:"user"`
	} `json:"graphql"`
	ToastContentOnLoad      interface{} `json:"toast_content_on_load"`
	ShowViewShop            bool        `json:"show_view_shop"`
	ProfilePicEditSyncProps interface{} `json:"profile_pic_edit_sync_props"`
}

func main() {
	url := fmt.Sprintf("https://www.instagram.com/pranav__satheesh/?__a=1")
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)

	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	resultStruct := InstagraMResponse{}
	body, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &resultStruct)
	log.Println(err)
	r := gin.Default()
	r.GET("/instahandle/:handle", func(c *gin.Context) {
		log.Println(c.Param("handle"))
		url := fmt.Sprintf("https://www.instagram.com/%s/?__a=1", c.Param("handle"))
		method := "GET"

		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"followers_count": 0,
				"result":          "failure",
			})
		}
		// req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/534.30 (KHTML, like Gecko) Ubuntu/20.04.1 Chromium/12.0.742.112 Chrome/12.0.742.112 Safari/534.30")
		req.Header.Set("User-Agent", "Golang_Spider_Bot/3.0")

		res, err := client.Do(req)

		if err != nil {
			fmt.Println(err)
			c.JSON(200, gin.H{
				"followers_count": 0,
				"result":          "failure",
			})
		}
		defer res.Body.Close()

		resultStruct := InstagraMResponse{}
		body, err := ioutil.ReadAll(res.Body)

		if err != nil {

			c.JSON(200, gin.H{
				"followers_count": 0,
				"result":          "failure",
			})
		}
		// log.Println(string(body))
		err = json.Unmarshal(body, &resultStruct)

		if err != nil {
			c.JSON(200, gin.H{
				"followers_count": 0,
				"result":          "failure",
			})
		}
		c.JSON(200, gin.H{
			"followers_count": resultStruct.Graphql.User.EdgeFollowedBy.Count,
			"result":          "success",
		})
	})
	r.Run(":3000")
}
