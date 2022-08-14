package controllers

import (
	"OnlineBooks/common"
	"OnlineBooks/models"
	"OnlineBooks/utils"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"html/template"
)

type BookController struct {
	BaseController
}

//我的图书页面
func (c *BookController) Index() {
	pageIndex, _ := c.GetInt("page", 1)
	private, _ := c.GetInt("private", 1) //默认私有
	books, totalCount, err := models.NewBook().SelectPage(pageIndex, common.PageSize, c.Member.MemberId, private)
	if err != nil {
		logs.Error("BookController.Index => ", err)
		c.Abort("404")
	}
	if totalCount > 0 {
		c.Data["PageHtml"] = utils.NewPaginations(common.RollPage, totalCount, common.PageSize, pageIndex, beego.URLFor("BookController.Index"), fmt.Sprintf("&private=%v", private))
	} else {
		c.Data["PageHtml"] = ""
	}
	//封面图片
	for idx, book := range books {
		book.Cover = utils.ShowImg(book.Cover, "cover")
		books[idx] = book
	}
	b, err := json.Marshal(books)
	if err != nil || len(books) <= 0 {
		c.Data["Result"] = template.JS("[]")
	} else {
		c.Data["Result"] = template.JS(string(b))
	}

	c.Data["Private"] = private
	c.TplName = "book/index.html"
}