package models

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

func init() {
	orm.RegisterModel(
		new(Category),
		new(Book),
		new(BookCategory),
		new(Document),
		new(DocumentStore),
		new(Attachment),
		new(Member),
		new(Relationship),
		new(Score),)
}

/*
* Table Names
*/

func TNCategory() string {
	return "category"
}

func TNBookCategory() string {
	return "book_category"
}

func TNBook() string {
	return "books"
}

func TNMembers() string {
	return "md_members"
}

func TNRelationship() string {
	return "md_relationship"
}

func TNDocuments() string {
	return "md_documents"
}

func TNAttachment() string {
	return "md_attachment"
}

func TNComments(bookid int) string {
	return fmt.Sprintf("md_comments_%04d", bookid%2)
}

func TNScore() string {
	return "md_score"
}

func TNDocumentStore() string {
	return "md_document_store"
}

//设置增减
//@param            table           需要处理的数据表
//@param            field           字段
//@param            condition       条件
//@param            incre           是否是增长值，true则增加，false则减少
//@param            step            增或减的步长
func IncOrDec(table string, field string, condition string, incre bool, step ...int) (err error) {
	mark := "-"
	if incre {
		mark = "+"
	}
	s := 1
	if len(step) > 0 {
		s = step[0]
	}
	sql := fmt.Sprintf("update %v set %v=%v%v%v where %v", table, field, field, mark, s, condition)
	_, err = orm.NewOrm().Raw(sql).Exec()
	return
}