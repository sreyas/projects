package controllers

import (
	"connect-staffing-inc/app/models"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/mattn/go-sqlite3"
	r "github.com/revel/revel"
)

var DB *gorm.DB

func init() {
	fmt.Println("in DB")
	var err error
	DB, err = gorm.Open("sqlite3", "../../sampleuser.db")
	if err != nil {
		panic(err.Error)
	} else {
		fmt.Println("Connected successfully")
	}
	initializeDb()
}

type ModelController struct {
	*r.Controller
	Orm *gorm.DB
}

func initializeDb() {
	migrateTable(&models.User{})
}
func migrateTable(m interface{}) {
	ex := DB.HasTable(m)
	if !ex {
		fmt.Println("migrating")
		if err := DB.AutoMigrate(m).Error; err != nil {
			fmt.Println("migrate erro ", err)
		}
	}
}
func (c *ModelController) Begin() r.Result {
	c.Orm = DB
	return nil
}