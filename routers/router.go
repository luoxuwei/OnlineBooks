package routers

import (
	"OnlineBooks/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.HomeController{}, "get:Index")
	beego.Router("/explore", &controllers.ExploreController{}, "get:Index")
	beego.Router("/books/:key", &controllers.DocumentController{}, "*:Index")

	//读书
	beego.Router("/read/:key/:id", &controllers.DocumentController{}, "*:Read")

	//编辑
	beego.Router("/api/:key/edit/?:id", &controllers.DocumentController{}, "*:Edit")
	beego.Router("/api/:key/content/?:id", &controllers.DocumentController{}, "*:Content")
	beego.Router("/api/upload", &controllers.DocumentController{}, "post:Upload")
	beego.Router("/api/:key/create", &controllers.DocumentController{}, "post:Create")
	beego.Router("/api/:key/delete", &controllers.DocumentController{}, "post:Delete")

	//图书管理
	beego.Router("/book", &controllers.BookController{}, "*:Index") //我的图书
	beego.Router("/book/create", &controllers.BookController{}, "post:Create") //创建图书
}
