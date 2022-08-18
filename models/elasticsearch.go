package models

import (
	"OnlineBooks/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"strconv"
	"strings"
)

func ElasticBuildIndex(bookId int) {
	//图书索引
	book, _ := NewBook().Select("book_id", bookId, "book_id", "book_name", "description")
	addBookToIndex(book.BookId, book.BookName, book.Description)

	//文档内容索引
	var documents []Document
	fields := []string{"document_id", "book_id", "document_name", "release"}
	GetOrm("r").QueryTable(TNDocuments()).Filter("book_id", bookId).All(&documents, fields...)
	if len(documents) > 0 {
		for _, document := range documents {
			addDocumentToIndex(document.DocumentId, document.BookId, flatHtml(document.Release))
		}
	}
}

func addBookToIndex(bookId int, bookName string, description string) {
	queryJson := `
		{
			"book_id":%v,
			"book_name":"%v",
			"description":"%v"
		}
	`

	//elasticsearch api
	host, _ := beego.AppConfig.String("elastic_host")
	api := host + "mbooks/datas/" + strconv.Itoa(bookId)

	//发起请求
	queryJson = fmt.Sprintf(queryJson, bookId, bookName, description)
	err := utils.HttpPutJson(api, queryJson)
	if nil != err {
		logs.Debug(err)
	}
}

func addDocumentToIndex(documentId, bookId int, release string) {
	queryJson := `
		{
			"document_id":%v,
			"book_id":%v,
			"release":"%v"
		}
	`

	//elasticsearch api
	host, _ := beego.AppConfig.String("elastic_host")
	api := host + "mdocuments/datas/" + strconv.Itoa(documentId)

	//发起请求
	queryJson = fmt.Sprintf(queryJson, documentId, bookId, release)
	err := utils.HttpPutJson(api, queryJson)
	if nil != err {
		logs.Debug(err)
	}

}

//一些像html标签的的内容与搜索无关，过滤掉
func flatHtml(htmlStr string) string {
	htmlStr = strings.Replace(htmlStr, "\n", " ", -1)
	htmlStr = strings.Replace(htmlStr, "\"", "", -1)

	gq, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr
	}
	return gq.Text()
}