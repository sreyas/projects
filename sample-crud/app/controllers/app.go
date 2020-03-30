package controllers

import (
	"connect-staffing-inc/app/models"
	"connect-staffing-inc/app/routes"
	"fmt"
	"strings"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	ModelController
}

func (c App) AddUser() revel.Result {
	if user := c.connected(); user != nil {
		c.ViewArgs["user"] = user
	}
	return nil
}

func (c App) connected() *models.User {
	if s, ok := c.Session["user"].(string); ok {
		c.cekSessionUser(s)
	}
	if username, ok := c.Session["user"].(string); ok {
		return c.getUser(username)
	}

	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}

	return nil
}

func (c App) cekSessionUser(username string) {
	var user models.User
	err := c.Orm.Where("username=?", username).Find(&user)
	if err.RecordNotFound() {
		c.Logout()
	}
}

func (c App) getUser(username string) *models.User {
	var user models.User
	err := c.Orm.Where("username=?", username).Find(&user)
	if err.RecordNotFound() {
		fmt.Println("Not Found")
	}
	return &user
}

func (c App) Index() revel.Result {
	if c.connected() != nil {
		return c.Redirect(routes.DashBoard.Index())
	}
	page := "login"
	usertype := ""
	return c.Render(usertype, page)

}

func (c App) Register() revel.Result {
	page := "login"
	usertype := ""
	return c.Render(page, usertype)
}

func (c App) Save(user models.User, password string, cekpassword string) revel.Result {
	c.Validation.Required(password).Message("Required pwd broooh")
	c.Validation.MaxSize(password, 255)
	c.Validation.MinSize(password, 7)
	c.Validation.Required(cekpassword).Message("Required Brooh")
	c.Validation.Required(cekpassword == password).Message("Not same")
	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.App.Register())
	}
	user.Password, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err := c.Orm.Save(&user)
	if err.Error != nil {
		c.Flash.Error("The username " + user.Username + " already exist")
		return c.Redirect(routes.App.Register())
	}
	c.Session["user"] = user.Username
	c.Flash.Success("welcome, " + strings.Title(user.Username))
	return c.Redirect(routes.PermPlace.Index())

}

func (c App) Login(username string, password string, remember bool) revel.Result {
	user := c.getUser(username)
	if user.Userstatus == "0" {
		c.Flash.Out["username"] = username
		c.Flash.Error("Account Disabled. Please contact Administrator")
		return c.Redirect(routes.App.Index())
	} else {
		if user != nil {
			err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
			if err == nil {
				c.Session["user"] = username
				c.Session["usertype"] = user.Usertype
				if remember {
					c.Session.SetDefaultExpiration()
				} else {
					c.Session.SetNoExpiration()
				}
				c.Flash.Success("welcome " + strings.Title(username))
				// return c.Redirect(routes.Users.Index())
				return c.Redirect(routes.DashBoard.Index())

			}
		}
	}
	c.Flash.Out["username"] = username
	c.Flash.Error("Login Failed")
	return c.Redirect(routes.App.Index())
}

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}
