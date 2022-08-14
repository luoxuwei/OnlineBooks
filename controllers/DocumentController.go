package controllers

import (
	"OnlineBooks/models"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"strings"
)

type DocumentController struct {
	BaseController
}

//获取图书内容并判断权限
func (c *DocumentController) getBookData(identify, token string) *models.BookData {
	book, err := models.NewBook().Select("identify", identify)
	if err != nil {
		logs.Error(err)
		c.Abort("404")
	}

	//私有文档
	if book.PrivatelyOwned == 1 && !c.Member.IsAdministrator() {
		isOk := false
		if c.Member != nil {
			_, err := models.NewRelationship().SelectRoleId(book.BookId, c.Member.MemberId)
			if err == nil {
				isOk = true
			}
		}
		if book.PrivateToken != "" && !isOk {
			if token != "" && strings.EqualFold(token, book.PrivateToken) {
				c.SetSession(identify, token)
			} else if token, ok := c.GetSession(identify).(string); !ok || !strings.EqualFold(token, book.PrivateToken) {
				c.Abort("404")
			}
		} else if !isOk {
			c.Abort("404")
		}
	}

	bookResult := book.ToBookData()
	if c.Member != nil {
		rsh, err := models.NewRelationship().Select(bookResult.BookId, c.Member.MemberId)
		if err == nil {
			bookResult.MemberId = rsh.MemberId
			bookResult.RoleId = rsh.RoleId
			bookResult.RelationshipId = rsh.RelationshipId
		}
	}
	return bookResult
}

//图书目录&详情页
func (c *DocumentController) Index() {
	token := c.GetString("token")
	identify := c.Ctx.Input.Param(":key")
	if identify == "" {
		c.Abort("404")
	}
	tab := strings.ToLower(c.GetString("tab"))

	bookResult := c.getBookData(identify, token)
	if bookResult.BookId == 0 { //没有阅读权限
		c.Redirect(beego.URLFor("HomeController.Index"), 302)
		return
	}

	c.TplName = "document/intro.html"
	c.Data["Book"] = bookResult

	switch tab {
	case "comment", "score":
	default:
		tab = "default"
	}
	c.Data["Tab"] = tab
	c.Data["Menu"], _ = new(models.Document).GetMenuTop(bookResult.BookId)

	c.Data["Comments"], _ = new(models.Comments).BookComments(1, 30, bookResult.BookId)
	c.Data["MyScore"] = new(models.Score).BookScoreByUid(c.Member.MemberId, bookResult.BookId)
}