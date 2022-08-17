package models

import (
	"github.com/beego/beego/v2/client/orm"
	"strconv"
)

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
	o := GetOrm("r")
	sql := "select c.* from " + TNCategory() + " c left join " + TNBookCategory() + " bc on c.id=bc.category_id where bc.book_id=?"
	rows, err = o.Raw(sql, book_id).QueryRows(&cates)
	return
}

//处理书籍分类
func (m *BookCategory) SetBookCates(bookId int, cids []string) {
	if len(cids) == 0 {
		return
	}

	var (
		cates             []Category
		tableCategory     = TNCategory()
		tableBookCategory = TNBookCategory()
	)

	o := GetOrm("w")
	o.QueryTable(tableCategory).Filter("id__in", cids).All(&cates, "id", "pid")

	cidMap := make(map[string]bool)
	for _, cate := range cates {
		cidMap[strconv.Itoa(cate.Pid)] = true
		cidMap[strconv.Itoa(cate.Id)] = true
	}
	cids = []string{}
	for cid, _ := range cidMap {
		cids = append(cids, cid)
	}

	o.QueryTable(tableBookCategory).Filter("book_id", bookId).Delete()
	var bookCates []BookCategory
	for _, cid := range cids {
		cidNum, _ := strconv.Atoi(cid)
		bookCate := BookCategory{
			CategoryId: cidNum,
			BookId:     bookId,
		}
		bookCates = append(bookCates, bookCate)
	}
	if l := len(bookCates); l > 0 {
		o.InsertMulti(l, &bookCates)
	}
	go CountCategory()
}

// 统计分类书籍
var counting = false
type Count struct {
	Cnt        int
	CategoryId int
}

func CountCategory() {
	if counting {
		return
	}
	counting = true
	defer func() {
		counting = false
	}()

	var count []Count

	o := GetOrm("w")
	sql := "select count(bc.id) cnt, bc.category_id from " + TNBookCategory() + " bc left join " + TNBook() + " b on b.book_id=bc.book_id where b.privately_owned=0 group by bc.category_id"
	o.Raw(sql).QueryRows(&count)
	if len(count) == 0 {
		return
	}

	var cates []Category
	o.QueryTable(TNCategory()).All(&cates, "id", "pid", "cnt")
	if len(cates) == 0 {
		return
	}

	var err error

	to, err := o.Begin()
	defer func() {
		if err != nil {
			to.Rollback()
		} else {
			to.Commit()
		}
	}()

	o.QueryTable(TNCategory()).Update(orm.Params{"cnt": 0})
	cateChild := make(map[int]int)
	for _, item := range count {
		if item.Cnt > 0 {
			cateChild[item.CategoryId] = item.Cnt
			_, err = o.QueryTable(TNCategory()).Filter("id", item.CategoryId).Update(orm.Params{"cnt": item.Cnt})
			if err != nil {
				return
			}
		}
	}
}
