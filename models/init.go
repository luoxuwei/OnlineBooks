package models

import "github.com/beego/beego/v2/client/orm"

func init() {
	orm.RegisterModel(
		new(Category),
		new(Book),
		new(BookCategory))
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