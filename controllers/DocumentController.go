package controllers

import (
	"OnlineBooks/common"
	"OnlineBooks/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"html/template"
	"strconv"
	"strings"
	"time"
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

//阅读器页面
func (c *DocumentController) Read() {
	identify := c.Ctx.Input.Param(":key")
	id := c.GetString(":id")
	token := c.GetString("token")

	if identify == "" || id == "" {
		c.Abort("404")
	}

	//没开启匿名
	if !c.EnableAnonymous && c.Member == nil {
		c.Redirect(beego.URLFor("AccountController.Login"), 302)
		return
	}

	bookData := c.getBookData(identify, token)

	doc := models.NewDocument()
	doc, err := doc.SelectByIdentify(bookData.BookId, id) //文档标识
	if err != nil {
		c.Abort("404")
	}

	if doc.BookId != bookData.BookId {
		c.Abort("404")
	}

	if doc.Release != "" {
		query, err := goquery.NewDocumentFromReader(bytes.NewBufferString(doc.Release))
		if err != nil {
			logs.Error(err)
		} else {
			query.Find("img").Each(func(i int, contentSelection *goquery.Selection) {
				if _, ok := contentSelection.Attr("src"); ok {
				}
				if alt, _ := contentSelection.Attr("alt"); alt == "" {
					contentSelection.SetAttr("alt", doc.DocumentName+" - 图"+fmt.Sprint(i+1))
				}
			})
			html, err := query.Find("body").Html()
			if err != nil {
				logs.Error(err)
			} else {
				doc.Release = html
			}
		}
	}

	attach, err := models.NewAttachment().SelectByDocumentId(doc.DocumentId)
	if err == nil {
		doc.AttachList = attach
	}

	//图书阅读人次+1
	if err := models.IncOrDec(models.TNBook(), "vcnt",
		fmt.Sprintf("book_id=%v", doc.BookId),
		true, 1,
	); err != nil {
		logs.Error(err.Error())
	}

	//文档阅读人次+1
	if err := models.IncOrDec(models.TNDocuments(), "vcnt",
		fmt.Sprintf("document_id=%v", doc.DocumentId),
		true, 1,
	); err != nil {
		logs.Error(err.Error())
	}

	doc.Vcnt = doc.Vcnt + 1

	if c.IsAjax() {
		var data struct {
			Id        int    `json:"doc_id"`
			DocTitle  string `json:"doc_title"`
			Body      string `json:"body"`
			Title     string `json:"title"`
			View      int    `json:"view"`
			UpdatedAt string `json:"updated_at"`
		}
		data.DocTitle = doc.DocumentName
		data.Body = doc.Release
		data.Id = doc.DocumentId
		data.View = doc.Vcnt
		data.UpdatedAt = doc.ModifyTime.Format("2006-01-02 15:04:05")

		c.JsonResult(0, "ok", data)
	}

	tree, err := models.NewDocument().GetMenuHtml(bookData.BookId, doc.DocumentId)
	if err != nil {
		logs.Error(err)
		c.Abort("404")
	}

	c.Data["Bookmark"] = false
	c.Data["Model"] = bookData
	c.Data["Book"] = bookData
	c.Data["Result"] = template.HTML(tree)
	c.Data["Title"] = doc.DocumentName
	c.Data["DocId"] = doc.DocumentId
	c.Data["Content"] = template.HTML(doc.Release)
	c.Data["View"] = doc.Vcnt
	c.Data["UpdatedAt"] = doc.ModifyTime.Format("2006-01-02 15:04:05")

	//设置模版
	c.TplName = "document/default_read.html"
}

//编辑
func (c *DocumentController) Edit() {
	docId := 0 // 文档id

	identify := c.Ctx.Input.Param(":key")
	if identify == "" {
		c.Abort("404")
	}

	bookData := models.NewBookData()

	var err error
	//权限验证
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().Select("identify", identify)
		if err != nil {
			c.JsonResult(1, "权限错误")
		}
		bookData = book.ToBookData()
	} else {
		bookData, err = models.NewBookData().SelectByIdentify(identify, c.Member.MemberId)
		if err != nil {
			c.Abort("404")
		}

		if bookData.RoleId == common.BookGeneral {
			c.JsonResult(1, "权限错误")
		}
	}

	c.TplName = "document/markdown_edit_template.html"

	c.Data["Model"] = bookData
	r, _ := json.Marshal(bookData)

	c.Data["ModelResult"] = template.JS(string(r))

	c.Data["Result"] = template.JS("[]")

	// 编辑的文档
	if id := c.GetString(":id"); id != "" {
		if num, _ := strconv.Atoi(id); num > 0 {
			docId = num
		} else { //字符串
			var doc = models.NewDocument()
			orm.NewOrm().QueryTable(doc).Filter("identify", id).Filter("book_id", bookData.BookId).One(doc, "document_id")
			docId = doc.DocumentId
		}
	}

	trees, err := models.NewDocument().GetMenu(bookData.BookId, docId, true)
	if err != nil {
		logs.Error("GetMenu error : ", err)
	} else {
		if len(trees) > 0 {
			if jsTree, err := json.Marshal(trees); err == nil {
				c.Data["Result"] = template.JS(string(jsTree))
			}
		} else {
			c.Data["Result"] = template.JS("[]")
		}
	}
	c.Data["BaiDuMapKey"] = beego.AppConfig.DefaultString("baidumapkey", "")

}

//保存文档并返回内容
func (c *DocumentController) Content() {
	identify := c.Ctx.Input.Param(":key")
	docId, err := c.GetInt("doc_id")
	errMsg := "ok"
	if err != nil {
		docId, _ = strconv.Atoi(c.Ctx.Input.Param(":id"))
	}
	bookId := 0
	//权限验证
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().Select("identify", identify)
		if err != nil {
			c.JsonResult(1, "获取内容错误")
		}
		bookId = book.BookId
	} else {
		bookData, err := models.NewBookData().SelectByIdentify(identify, c.Member.MemberId)

		if err != nil || bookData.RoleId == common.BookGeneral {
			c.JsonResult(1, "权限错误")
		}
		bookId = bookData.BookId
	}

	if docId <= 0 {
		c.JsonResult(1, "参数错误")
	}

	documentStore := new(models.DocumentStore)

	if !c.Ctx.Input.IsPost() {
		doc, err := models.NewDocument().SelectByDocId(docId)

		if err != nil {
			c.JsonResult(1, "文档不存在")
		}
		attach, err := models.NewAttachment().SelectByDocumentId(doc.DocumentId)
		if err == nil {
			doc.AttachList = attach
		}

		doc.Release = "" //Ajax请求，之间用markdown渲染，不用release
		doc.Markdown = documentStore.SelectField(doc.DocumentId, "markdown")
		c.JsonResult(0, errMsg, doc)
	}

	//更新文档内容
	markdown := strings.TrimSpace(c.GetString("markdown", ""))
	content := c.GetString("html")

	version, _ := c.GetInt64("version", 0)
	isCover := c.GetString("cover")

	doc, err := models.NewDocument().SelectByDocId(docId)

	if err != nil {
		c.JsonResult(1, "读取文档错误")
	}
	if doc.BookId != bookId {
		c.JsonResult(1, "内部错误")
	}
	if doc.Version != version && !strings.EqualFold(isCover, "yes") {
		c.JsonResult(1, "文档将被覆盖")
	}

	isSummary := false
	isAuto := false

	if markdown == "" && content != "" {
		documentStore.Markdown = content
	} else {
		documentStore.Markdown = markdown
	}
	documentStore.Content = content
	doc.Version = time.Now().Unix()
	if docId, err := doc.InsertOrUpdate(); err != nil {
		c.JsonResult(1, "保存失败")
	} else {
		documentStore.DocumentId = int(docId)
		if err := documentStore.InsertOrUpdate("markdown", "content"); err != nil {
			logs.Error(err)
		}
	}

	if isAuto {
		errMsg = "auto"
	} else if isSummary {
		errMsg = "true"
	}

	doc.Release = ""
	c.JsonResult(0, errMsg, doc)
}