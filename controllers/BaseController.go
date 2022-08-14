package controllers

import (
	"OnlineBooks/models"
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
	Member          *models.Member    //用户
}