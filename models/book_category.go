package models

import "github.com/beego/beego/v2/client/orm"

//图书分类对应关系
type BookCategory struct {
	Id         int //自增主键
	BookId     int //书籍id
	CategoryId int //分类id
}

func (m *BookCategory) TableName() string {
	return TNBookCategory()
}


//根据书籍id查询分类id
func (m *BookCategory) SelectByBookId(book_id int) (cates []Category, rows int64, err error) {
	o := orm.NewOrm()
	sql := "select c.* from " + TNCategory() + " c left join " + TNBookCategory() + " bc on c.id=bc.category_id where bc.book_id=?"
	rows, err = o.Raw(sql, book_id).QueryRows(&cates)
	return
}

