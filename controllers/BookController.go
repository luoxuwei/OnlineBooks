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
	"strings"
	"time"
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

//创建图书
func (c *BookController) Create() {
	identify := strings.TrimSpace(c.GetString("identify", ""))
	bookName := strings.TrimSpace(c.GetString("book_name", ""))
	author := strings.TrimSpace(c.GetString("author", ""))
	authorURL := strings.TrimSpace(c.GetString("author_url", ""))
	privatelyOwned, _ := c.GetInt("privately_owned")
	description := strings.TrimSpace(c.GetString("description", ""))

	/*
	* 约束条件判断
	 */
	if identify == "" || strings.Count(identify, "") > 50 {
		c.JsonResult(1, "请正确填写图书标识，不能超过50字")
	}
	if bookName == "" {
		c.JsonResult(1, "请填图书名称")
	}

	if strings.Count(description, "") > 500 {
		c.JsonResult(1, "图书描述需小于500字")
	}

	if privatelyOwned != 0 && privatelyOwned != 1 {
		privatelyOwned = 1
	}

	book := models.NewBook()
	if book, _ := book.Select("identify", identify); book.BookId > 0 {
		c.JsonResult(1, "identify冲突")
	}

	book.BookName = bookName
	book.Identify = identify
	book.Description = description
	book.CommentCount = 0
	book.PrivatelyOwned = privatelyOwned
	book.Cover = common.DefaultCover()
	book.DocCount = 0
	book.MemberId = c.Member.MemberId
	book.CommentCount = 0
	book.Editor = "markdown"
	book.ReleaseTime = time.Now()
	book.Score = 40 //评分
	book.Author = author
	book.AuthorURL = authorURL

	if err := book.Insert(); err != nil {
		c.JsonResult(1, "数据库错误")
	}

	bookResult, err := models.NewBookData().SelectByIdentify(book.Identify, c.Member.MemberId)
	if err != nil {
		logs.Error(err)
	}

	c.JsonResult(0, "ok", bookResult)
}