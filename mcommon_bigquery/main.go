package main

import (
	"mcommon_bigquery/controllers"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
)

// Mcommons Details Insert into Google Bigquery
func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping",
		})
	})

	r.GET("/bigquery", func(c *gin.Context) {

		go controllers.GetBroadCastDetails()
		go controllers.GetProfileDetails()
		go controllers.GetMessages()
		go controllers.GetSentMessages()
		go controllers.GetKeyWords()

		go controllers.GetTinyURLS()
		go controllers.GetClicksDetails()

		c.JSON(200, gin.H{
			"status": "true",
		})
	})

	r.Run(":3000")

}

type Item struct {
	Name string
	Age  string
}

func (i *Item) Save() (map[string]bigquery.Value, string, error) {
	return map[string]bigquery.Value{
		"name": i.Name,
		"age":  i.Age,
	}, bigquery.NoDedupeID, nil
}
