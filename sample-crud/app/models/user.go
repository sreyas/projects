package models

import (
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
	"regexp"
)

type User struct {
	gorm.Model
	Username   string `gorm:"size:100;unique"`
	Password   []byte
	Email      string `gorm:"size:100"`
	Firstname  string `gorm:"size:100"`
	Lastname   string `gorm:"size:100"`
	Usertype   string `gorm:"size:1"`
	Userstatus string `gorm:"size:1"`
}

var userRegex = regexp.MustCompile("^\\w*$")

func (user *User) Validate(v *revel.Validation) {
	v.Check(
		user.Firstname,
		revel.Required{},
		revel.MinSize{2},
		revel.MaxSize{100},
	)

	v.Check(
		user.Username,
		revel.Required{},
		revel.MinSize{4},
		revel.MaxSize{100},
		revel.Match{userRegex},
	)
}
