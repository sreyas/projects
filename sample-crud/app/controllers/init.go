package controllers

import "github.com/revel/revel"

func init() {
	revel.InterceptMethod((*ModelController).Begin, revel.BEFORE)
	revel.InterceptMethod(App.AddUser, revel.BEFORE)
	revel.InterceptMethod(Users.cekUser, revel.BEFORE)
}
