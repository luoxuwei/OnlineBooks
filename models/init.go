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