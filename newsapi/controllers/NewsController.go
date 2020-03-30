package controllers

import (
	"app/newsapi/newsapi/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

const baseurl = "https://newsapi.org/v2/everything"

type NewsController struct {
	BaseController
}
type Data struct {
	Name string
}
type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client
	Options    url.Values
}
type Resp struct {
	Response *models.Response
}

type Options struct {
	Q        string
	SortBy   string
	APIKey   string
	pageSize string
	Language string
}

/*
Loads the search page


*/
func (n *NewsController) Index() {
	n.TplName = "template/index.html"

}

/*
NewsRead returns news corresponding to query
*/
func (n *NewsController) NewsRead() {
	w := n.Ctx.ResponseWriter
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,/")
	query := n.Ctx.Input.Param(":data") //Takes query from request
	query = strings.ToLower(query)      //Convert it into lower character
	crntDate := time.Now().Local().Format("2006-01-02 15:04:05")
	queryStatus := models.CheckQueryExist(query, crntDate) //Checks in Db whether the query was the previous one
	news := &models.Response{}
	if queryStatus > 0 {
		news = models.GetNewsFromDB(query) // If the query was from previous then takes news details which were saved in db.
	} else {

		client := NewClient()
		apiKey := beego.AppConfig.String("api_key") //API Key retrieves from conf file
		options := Options{Q: query, SortBy: "publishedAt", APIKey: apiKey, pageSize: "10", Language: "en"}
		client.Options = options.AddOptions()
		result, err := client.GetNews()
		news = result.Response
		log.Println(err)
		models.InsertNewsContent(news, query, crntDate)
	}
	slice := []interface{}{*news}
	sliceToClient, _ := json.Marshal(slice[0])
	w.Write(sliceToClient)
}

/*
NewClient returns a Client which has
	A http client
	Parsed Url
*/
func NewClient() *Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	httpClient := http.Client{Transport: tr}
	client := Client{}
	client.httpClient = &httpClient
	client.BaseURL, _ = url.Parse(baseurl)
	return &client
}

/*
Addoptions return url.Values from Options Struct
*/
func (o *Options) AddOptions() url.Values {
	urlValues := url.Values{}
	urlValues.Add("q", o.Q)
	urlValues.Add("sortBy", o.SortBy)
	urlValues.Add("apiKey", o.APIKey)
	urlValues.Add("pageSize", o.pageSize)
	urlValues.Add("language", o.Language)
	return urlValues
}

/*
GetNews returns news details from online
It makes complete url using fmt.Sprintf -> (parsedurl+url options)
Makes a GET request and the response parse into a struct which is of the response format.
*/
func (c *Client) GetNews() (Resp, error) {
	url := fmt.Sprintf("%s?%s", c.BaseURL.String(), c.Options.Encode())
	req, err := http.NewRequest("GET", url, nil)
	re := Resp{}

	if err != nil {
		return re, err
	}
	response := new(models.Response)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return re, err
	}
	if resp.StatusCode != 200 {
		return re, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(response)
	re.Response = response
	return re, err
}
