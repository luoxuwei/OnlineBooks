package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"strings"
	"time"
)

//图书章节内容
type Document struct {
	DocumentId   int           `orm:"pk;auto;column(document_id)" json:"doc_id"`
	DocumentName string        `orm:"column(document_name);size(500)" json:"doc_name"`
	Identify     string        `orm:"column(identify);size(100);index;null;default(null)" json:"identify"`
	BookId       int           `orm:"column(book_id);type(int)" json:"book_id"`
	ParentId     int           `orm:"column(parent_id);type(int);default(0)" json:"parent_id"`
	OrderSort    int           `orm:"column(order_sort);default(0);type(int)" json:"order_sort"`
	Release      string        `orm:"column(release);type(text);null" json:"release"`
	CreateTime   time.Time     `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	MemberId     int           `orm:"column(member_id);type(int)" json:"member_id"`
	ModifyTime   time.Time     `orm:"column(modify_time);type(datetime);default(null);auto_now" json:"modify_time"`
	ModifyAt     int           `orm:"column(modify_at);type(int)" json:"-"`
	Version      int64         `orm:"type(bigint);column(version)" json:"version"`
	AttachList   []*Attachment `orm:"-" json:"attach"`
	Vcnt         int           `orm:"column(vcnt);default(0)" json:"vcnt"`
	Markdown     string        `orm:"-" json:"markdown"`
}

func (m *Document) TableName() string {
	return TNDocuments()
}

func NewDocument() *Document {
	return &Document{
		Version: time.Now().Unix(),
	}
}


func (m *Document) GetMenuTop(bookId int) (docs []*Document, err error) {
	var docsAll []*Document
	o := orm.NewOrm()
	cols := []string{"document_id", "document_name", "member_id", "parent_id", "book_id", "identify"}
	fmt.Println("---------------start")
	_, err = o.QueryTable(m.TableName()).Filter("book_id", bookId).Filter("parent_id", 0).OrderBy("order_sort", "document_id").Limit(5000).All(&docsAll, cols...)
	fmt.Println("---------------end")
	for _, doc := range docsAll {
		docs = append(docs, doc)
	}
	return
}

//根据指定字段查询一条文档
func (m *Document) SelectByIdentify(BookId, Identify interface{}) (*Document, error) {
	err := orm.NewOrm().QueryTable(m.TableName()).Filter("BookId", BookId).Filter("Identify", Identify).One(m)
	return m, err
}

//根据文档ID查询指定文档
func (m *Document) SelectByDocId(id int) (doc *Document, err error) {
	if id <= 0 {
		return m, errors.New("Invalid parameter")
	}

	o := orm.NewOrm()
	err = o.QueryTable(m.TableName()).Filter("document_id", id).One(m)
	if err == orm.ErrNoRows {
		return m, errors.New("数据不存在")
	}

	return m, nil
}

//插入和更新文档
func (m *Document) InsertOrUpdate(cols ...string) (id int64, err error) {
	o := orm.NewOrm()
	id = int64(m.DocumentId)
	m.ModifyTime = time.Now()
	m.DocumentName = strings.TrimSpace(m.DocumentName)
	if m.DocumentId > 0 { //文档id存在，则更新
		_, err = o.Update(m, cols...)
		return
	}

	var selectDocument Document
	//直接查询一个字段
	o.QueryTable(TNDocuments()).Filter("identify", m.Identify).Filter("book_id", m.BookId).One(&selectDocument, "document_id")
	if selectDocument.DocumentId == 0 {
		m.CreateTime = time.Now()
		id, err = o.Insert(m)
		NewBook().RefreshDocumentCount(m.BookId)
	} else { //identify存在，则执行更新
		_, err = o.Update(m)
		id = int64(selectDocument.DocumentId)
	}
	return
}

//删除文档及其子文档
func (m *Document) Delete(docId int) error {

	o := orm.NewOrm()
	modelStore := new(DocumentStore)

	if doc, err := m.SelectByDocId(docId); err == nil {
		o.Delete(doc)
		modelStore.Delete(docId)
	}

	var docs []*Document

	_, err := o.QueryTable(m.TableName()).Filter("parent_id", docId).All(&docs)
	if err != nil {
		return err
	}

	for _, item := range docs {
		docId := item.DocumentId
		o.QueryTable(m.TableName()).Filter("document_id", docId).Delete()
		//删除document_store表对应的文档
		modelStore.Delete(docId)
		m.Delete(docId)
	}
	return nil
}
