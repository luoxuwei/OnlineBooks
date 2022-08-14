package controllers

import (
	"OnlineBooks/models"
	"OnlineBooks/utils"
	"github.com/beego/beego/v2/core/logs"
	"strings"
)

type ManagerController struct {
	BaseController
}

//分类管理
func (c *ManagerController) Category() {
	cate := new(models.Category)
	if strings.ToLower(c.Ctx.Request.Method) == "post" {
		//新增分类
		pid, _ := c.GetInt("pid")
		if err := cate.InsertMulti(pid, c.GetString("cates")); err != nil {
			c.JsonResult(1, "新增失败："+err.Error())
		}
		c.JsonResult(0, "新增成功")
	}

	//查询所有分类
	cates, err := cate.GetCates(-1, -1)
	if err != nil {
		logs.Error(err)
	}

	var parents []models.Category
	for idx, item := range cates {
		if strings.TrimSpace(item.Icon) == "" {
			item.Icon = "/static/images/icon.png"
		} else {
			item.Icon = utils.ShowImg(item.Icon)
		}
		if item.Pid == 0 {
			parents = append(parents, item)
		}
		cates[idx] = item
	}

	c.Data["Parents"] = parents
	c.Data["Cates"] = cates
	c.Data["IsCategory"] = true
	c.TplName = "manager/category.html"
}