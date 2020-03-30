package controllers

import (
	"connect-staffing-inc/app/models"
	"connect-staffing-inc/app/routes"
	"strings"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	App
}

func (c Users) cekUser() revel.Result {
	if user := c.connected(); user == nil {
		return c.Redirect(routes.App.Index())
	}
	return nil
}

func (c Users) Index() revel.Result {
	var user []models.User
	userobj := c.Orm.Find(&user)
	if userobj.Error != nil {
		panic(userobj.Error)
	}

	page := "user"
	usrDetails := c.connected()
	name := strings.Title(usrDetails.Firstname)
	usertype := usrDetails.Usertype
	if usrDetails.Usertype != "1" {
		return c.Redirect(routes.DashBoard.Index())
	}
	return c.Render(user, name, usertype, page)
}

func (c Users) New() revel.Result {
	page := "user"
	usrDetails := c.connected()
	name := strings.Title(usrDetails.Firstname)
	usertype := usrDetails.Usertype
	return c.Render(name, usertype, page)
}

func (c Users) Save(user models.User, password string, cekpassword string) revel.Result {
	c.Validation.Required(password).Message("Required Password")
	c.Validation.MaxSize(password, 255)
	c.Validation.MinSize(password, 7)
	c.Validation.Required(cekpassword).Message("Verify Password")
	c.Validation.Required(cekpassword == password).Message("Not same")
	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Users.New())
	}
	user.Password, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err := c.Orm.Save(&user)
	if err.Error != nil {
		c.Flash.Error("The username " + user.Username + " already exist")
		return c.Redirect(routes.Users.New())
	}
	return c.Redirect(routes.Users.Index())
}

func (c Users) Edit(id int) revel.Result {
	userobj := c.getUsers(id)

	usrDetails := c.connected()
	name := strings.Title(usrDetails.Firstname)
	usertype := usrDetails.Usertype
	page := "user"
	return c.Render(userobj, name, usertype, page)
}

func (c Users) Update(id int) revel.Result {

	userobj := c.getUsers(id)

	userobj.Firstname = c.Params.Get("Firstname")
	userobj.Lastname = c.Params.Get("Lastname")
	userobj.Email = c.Params.Get("Email")
	userobj.Usertype = c.Params.Get("Usertype")
	userobj.Userstatus = c.Params.Get("Userstatus")

	password := c.Params.Get("password")
	cekpassword := c.Params.Get("cekpassword")
	if password != "" || cekpassword != "" {
		c.Validation.Required(password).Message("Required Password")
		c.Validation.MaxSize(password, 255)
		c.Validation.MinSize(password, 7)
		c.Validation.Required(cekpassword).Message("Verify Password")
		c.Validation.Required(cekpassword == password).Message("Not same")
		userobj.Validate(c.Validation)
		if c.Validation.HasErrors() {
			c.Validation.Keep()
			c.FlashParams()
			return c.Redirect(routes.Users.Edit(id))
		}
		userobj.Password, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	}

	c.Validation.Required(userobj.Firstname).Message("Firstname cannot be empty")
	c.Validation.Required(userobj.Lastname).Message("Lastname cannot be empty")
	userobj.Validate(c.Validation)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Users.Edit(id))
	}

	err := c.Orm.Save(&userobj)
	if err.Error != nil {
		panic(err.Error)
	}
	return c.Redirect(routes.Users.Index())
}

func (c Users) Show(id int) revel.Result {
	userobj := c.getUsers(id)
	usrDetails := c.connected()
	name := strings.Title(usrDetails.Firstname)
	usertype := usrDetails.Usertype
	page := "user"
	return c.Render(userobj, name, usertype, page)
}

func (c Users) Delete(id int) revel.Result {
	var pp models.User
	err := c.Orm.Where("ID=?", id).Delete(&pp)
	if err.Error != nil {
		panic(err.Error)
	}
	return c.Redirect(routes.Users.Index())

}

func (c Users) getUsers(id int) *models.User {
	var user models.User
	userobj := c.Orm.Find(&user, id)
	if userobj.Error != nil {
		panic(userobj.Error)
	}
	return &user
}
