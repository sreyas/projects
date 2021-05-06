package main

import (
	"mcommon/htmllink"

	"github.com/gin-gonic/gin"
)

//Mcommon BroadCast Details

func main() {

	r := gin.Default()
	r.GET("/getBroadCast", func(c *gin.Context) {

		campaignId := c.DefaultQuery("campaignId", "")
		broadCastId := c.DefaultQuery("broadCastId", "")
		mcommnsSessionId := c.DefaultQuery("sessionId", "")
		broadcast := htmllink.ParseHtml(campaignId, broadCastId, mcommnsSessionId)
		c.JSON(200, gin.H{
			"status": "true",
			"result": broadcast,
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping",
		})
	})
	r.Run(":8080")
}
