package controllers

import (
	"OnlineBooks/common"
	"OnlineBooks/models"
	"OnlineBooks/utils"
	"OnlineBooks/utils/graphics"
	"OnlineBooks/utils/store"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
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

// 设置图书页面
func (c *BookController) Setting() {

	key := c.Ctx.Input.Param(":key")

	if key == "" {
		c.Abort("404")
	}

	book, err := models.NewBookData().SelectByIdentify(key, c.Member.MemberId)
	if err != nil && err != orm.ErrNoRows {
		c.Abort("404")
	}

	//需管理员以上权限
	if book.RoleId != common.BookFounder && book.RoleId != common.BookAdmin {
		c.Abort("404")
	}

	if book.PrivateToken != "" {
		book.PrivateToken = c.BaseUrl() + beego.URLFor("DocumentController.Index", ":key", book.Identify, "token", book.PrivateToken)
	}

	//查询图书分类
	if selectedCates, rows, _ := new(models.BookCategory).SelectByBookId(book.BookId); rows > 0 {
		var maps = make(map[int]bool)
		for _, cate := range selectedCates {
			maps[cate.Id] = true
		}
		c.Data["Maps"] = maps
	}

	c.Data["Cates"], _ = new(models.Category).GetCates(-1, 1)
	c.Data["Model"] = book
	c.TplName = "book/setting.html"
}

func (c *BookController) isPermission() (*models.BookData, error) {

	identify := c.GetString("identify")

	book, err := models.NewBookData().SelectByIdentify(identify, c.Member.MemberId)
	if err != nil {
		return book, err
	}

	if book.RoleId != common.BookAdmin && book.RoleId != common.BookFounder {
		return book, errors.New("权限不足")
	}
	return book, nil
}

//上传封面.
func (c *BookController) UploadCover() {
	bookResult, err := c.isPermission()
	if err != nil {
		c.JsonResult(1, err.Error())
	}

	book, err := models.NewBook().Select("book_id", bookResult.BookId)
	if err != nil {
		c.JsonResult(1, err.Error())
	}

	file, moreFile, err := c.GetFile("image-file")
	if err != nil {
		logs.Error("", err.Error())
		c.JsonResult(1, "读取文件异常")
	}

	defer file.Close()

	ext := filepath.Ext(moreFile.Filename)

	if !strings.EqualFold(ext, ".png") && !strings.EqualFold(ext, ".jpg") && !strings.EqualFold(ext, ".gif") && !strings.EqualFold(ext, ".jpeg") {
		c.JsonResult(1, "不支持图片格式")
	}

	x1, _ := strconv.ParseFloat(c.GetString("x"), 10)
	y1, _ := strconv.ParseFloat(c.GetString("y"), 10)
	w1, _ := strconv.ParseFloat(c.GetString("width"), 10)
	h1, _ := strconv.ParseFloat(c.GetString("height"), 10)

	x := int(x1)
	y := int(y1)
	width := int(w1)
	height := int(h1)

	fileName := strconv.FormatInt(time.Now().UnixNano(), 16)

	filePath := filepath.Join("uploads", time.Now().Format("200601"), fileName+ext)

	path := filepath.Dir(filePath)

	os.MkdirAll(path, os.ModePerm)

	err = c.SaveToFile("image-file", filePath)

	if err != nil {
		logs.Error("", err)
		c.JsonResult(1, "保存图片失败")
	}

	//剪切图片
	subImg, err := graphics.ImageCopyFromFile(filePath, x, y, width, height)
	if err != nil {
		c.JsonResult(1, "图片剪切")
	}

	filePath = filepath.Join(common.WorkingDirectory, "uploads", time.Now().Format("200601"), fileName+ext)

	//生成缩略图
	err = graphics.ImageResizeSaveFile(subImg, 175, 230, filePath)
	if err != nil {
		c.JsonResult(1, "保存图片失败")
	}

	url := "/" + strings.Replace(strings.TrimPrefix(filePath, common.WorkingDirectory), "\\", "/", -1)
	if strings.HasPrefix(url, "//") {
		url = string(url[1:])
	}
	book.Cover = url

	if err := book.Update(); err != nil {
		c.JsonResult(1, "保存图片失败")
	}

	save := book.Cover
	if err := store.SaveToLocal("."+url, save); err != nil {
		logs.Error(err.Error())
	} else {
		url = book.Cover
	}
	c.JsonResult(0, "ok", url)
}

//收藏
func (c *BookController) Collection() {
	uid := c.BaseController.Member.MemberId
	if uid <= 0 {
		c.JsonResult(1, "收藏失败，请先登录")
	}

	id, _ := c.GetInt(":id")
	if id <= 0 {
		c.JsonResult(1, "收藏失败，图书不存在")
	}

	cancel, err := new(models.Collection).Collection(uid, id)
	data := map[string]bool{"IsCancel": cancel}
	if err != nil {
		logs.Error(err.Error())
		if cancel {
			c.JsonResult(1, "取消收藏失败", data)
		}
		c.JsonResult(1, "添加收藏失败", data)
	}

	if cancel {
		c.JsonResult(0, "取消收藏成功", data)
	}
	c.JsonResult(0, "添加收藏成功", data)
}

//保存图书信息
func (c *BookController) SaveBook() {

	bookResult, err := c.isPermission()
	if err != nil {
		c.JsonResult(1, err.Error())
	}

	book, err := models.NewBook().Select("book_id", bookResult.BookId)
	if err != nil {
		logs.Error("SaveBook => ", err)
		c.JsonResult(1, err.Error())
	}

	bookName := strings.TrimSpace(c.GetString("book_name"))
	description := strings.TrimSpace(c.GetString("description", ""))
	editor := strings.TrimSpace(c.GetString("editor"))

	if strings.Count(description, "") > 500 {
		c.JsonResult(1, "描述需小于500字")
	}

	if editor != "markdown" && editor != "html" {
		editor = "markdown"
	}

	book.BookName = bookName
	book.Description = description
	book.Editor = editor
	book.Author = c.GetString("author")
	book.AuthorURL = c.GetString("author_url")

	if err := book.Update(); err != nil {
		c.JsonResult(1, "保存失败")
	}
	bookResult.BookName = bookName
	bookResult.Description = description

	//Update分类
	if cids, ok := c.Ctx.Request.Form["cid"]; ok {
		new(models.BookCategory).SetBookCates(book.BookId, cids)
	}

	c.JsonResult(0, "ok", bookResult)
}

//发布图书.
func (c *BookController) Release() {
	identify := c.GetString("identify")
	bookId := 0
	if c.Member.IsAdministrator() {
		book, err := models.NewBook().Select("identify", identify)
		if err != nil {
			logs.Error(err)
		}
		bookId = book.BookId
	} else {
		book, err := models.NewBookData().SelectByIdentify(identify, c.Member.MemberId)
		if err != nil {
			c.JsonResult(1, "未知错误")
		}
		if book.RoleId != common.BookAdmin && book.RoleId != common.BookFounder && book.RoleId != common.BookEditor {
			c.JsonResult(1, "权限不足")
		}
		bookId = book.BookId
	}

	if exist := utils.BooksRelease.Exist(bookId); exist {
		c.JsonResult(1, "正在发布中，请稍后操作")
	}

	go func(identify string) {
		models.NewDocument().ReleaseContent(bookId, c.BaseUrl())
	}(identify)

	c.JsonResult(0, "已发布")
}