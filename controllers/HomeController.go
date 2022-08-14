package controllers

import (
	"OnlineBooks/models"
	"github.com/beego/beego/v2/core/logs"
)

type HomeController struct {
	BaseController
}


func (c *HomeController) Index() {
	if cates, err := new(models.Category).GetCates(-1, 1); err == nil {
		c.Data["Cates"] = cates
	} else {
		logs.Error(err.Error())
	}

	c.TplName = "home/list.html"
}