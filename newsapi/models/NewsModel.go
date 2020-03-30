package models

import (
	"database/sql"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

type DBDetails struct {
	DbDriver   string
	DbUser     string
	Dbpassword string
	DbName     string
}
type Response struct {
	News
}
type News struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}
type Article struct {
	Source      Source
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	URLToImage  string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
}
type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

/*
dbConn opens a database connection.
It takes db driver, user name, password and db name from conf file
*/
func dbConn() *sql.DB {
	dbDet := DBDetails{}
	dbDet.DbDriver = beego.AppConfig.String("dbDriver")
	dbDet.DbUser = beego.AppConfig.String("dbUser")
	dbDet.Dbpassword = beego.AppConfig.String("dbPass")
	dbDet.DbName = beego.AppConfig.String("dbName")
	db, err := sql.Open(dbDet.DbDriver, dbDet.DbUser+":"+dbDet.Dbpassword+"@/"+dbDet.DbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

/*
CheckQueryExist returns 0 if there is the same query wasn't the last used or query was used a day ago.
Returns greater than 0 value if the query is the last one used on the same day.
*/
func CheckQueryExist(query string, crntDate string) int {
	db := dbConn()
	defer db.Close()
	qryExist := 0
	selStmt := db.QueryRow("SELECT COUNT(*) FROM query WHERE QueryName=? AND DATE(CreatedDate)=DATE(?)", query, crntDate)
	err := selStmt.Scan(&qryExist)
	if err != nil {
		panic(err.Error())
	}
	return qryExist
}

/*
GetNews Returns all news Details corresponding to the query

*/
func GetNewsFromDB(query string) *Response {
	db := dbConn()
	defer db.Close()
	resp := new(Response)
	resp.News = News{}
	resp.News.Status = "ok"
	resCount := 0
	resp.News.Articles = []Article{}
	SelNews, err := db.Query("SELECT  ContentAuthor, ContentDescription, ContentPublishedAt, ContentTitle, ContentUrl, ContentUrlImage, ContentSourceId, ContentSourceName FROM query LEFT JOIN content ON content.ContentQueryName = query.QueryName where query.QueryName=?", query)
	if err != nil {
		panic(err.Error())
	}
	for SelNews.Next() {
		resCount++
		articleData := Article{}
		sourceData := Source{}
		var authorName, description, publishedAt, contentUrl, urlToImage, sourceId, sourceName sql.NullString
		err := SelNews.Scan(&authorName, &description, &publishedAt, &articleData.Title, &contentUrl, &urlToImage, &sourceId, &sourceName)
		if err != nil {
			panic(err.Error())
		}
		articleData.Author = authorName.String
		articleData.Description = description.String
		articleData.PublishedAt = publishedAt.String
		articleData.URL = contentUrl.String
		articleData.URLToImage = urlToImage.String
		sourceData.ID = sourceId.String
		sourceData.Name = sourceName.String
		articleData.Source = sourceData
		resp.News.Articles = append(resp.News.Articles, articleData)
	}
	resp.News.TotalResults = resCount
	return resp
}
func InsertNewsContent(r *Response, query string, crntDate string) {
	db := dbConn()
	defer db.Close()
	res, err := db.Exec("truncate table query")
	res, err = db.Exec("truncate table content")
	tx, _ := db.Begin()
	instStmt, err := db.Prepare("INSERT INTO query (QueryName,CreatedDate) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	res, err = instStmt.Exec(query, crntDate)
	tx.Commit()
	affect, _ := res.RowsAffected()
	if affect > 0 {
		for _, x := range r.News.Articles {
			tx, _ = db.Begin()
			cntntInstStmt, err := tx.Prepare("INSERT INTO content (ContentQueryName, ContentAuthor, ContentDescription, ContentPublishedAt, ContentTitle, ContentUrl, ContentUrlImage, ContentSourceId, ContentSourceName) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")

			if err != nil {
				panic(err.Error())
			}
			_, err = cntntInstStmt.Exec(query, x.Author, x.Description, x.PublishedAt, x.Title, x.URL, x.URLToImage, x.Source.ID, x.Source.Name)
			tx.Commit()
		}
	}
}
