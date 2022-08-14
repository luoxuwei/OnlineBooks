package routers

import (
	"OnlineBooks/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.HomeController{}, "get:Index")
	beego.Router("/explore", &controllers.ExploreController{}, "get:Index")
	beego.Router("/books/:key", &controllers.DocumentController{}, "*:Index")
}
