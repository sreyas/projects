package main

import (
	"app/newsapi/newsapi/controllers"
	_ "app/newsapi/newsapi/routers"

	"github.com/astaxie/beego"
)

func main() {
	beego.Router("/index",&controllers.NewsController{},"*:Index")
	beego.Router("/news/:data", &controllers.NewsController{}, "*:NewsRead")
	beego.Run()
}
