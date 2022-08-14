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
	beego.Router("/read/:key/search", &controllers.DocumentController{}, "post:Search")

	//搜索
	beego.Router("/search", &controllers.SearchController{}, "get:Search")
	beego.Router("/search/result", &controllers.SearchController{}, "get:Result")

	//编辑
	beego.Router("/api/:key/edit/?:id", &controllers.DocumentController{}, "*:Edit")
	beego.Router("/api/:key/content/?:id", &controllers.DocumentController{}, "*:Content")
	beego.Router("/api/upload", &controllers.DocumentController{}, "post:Upload")
	beego.Router("/api/:key/create", &controllers.DocumentController{}, "post:Create")
	beego.Router("/api/:key/delete", &controllers.DocumentController{}, "post:Delete")

	//图书管理
	beego.Router("/book", &controllers.BookController{}, "*:Index") //我的图书
	beego.Router("/book/create", &controllers.BookController{}, "post:Create") //创建图书
	beego.Router("/book/:key/setting", &controllers.BookController{}, "*:Setting") //图书设置
	beego.Router("/book/setting/upload", &controllers.BookController{}, "post:UploadCover") //图书封面
	beego.Router("/book/star/:id", &controllers.BookController{}, "*:Collection")           //收藏图书
	beego.Router("/book/setting/save", &controllers.BookController{}, "post:SaveBook")      //保存
	beego.Router("/book/:key/release", &controllers.BookController{}, "post:Release")       //发布
	beego.Router("/book/setting/token", &controllers.BookController{}, "post:CreateToken")  //创建Token

	//管理后台
	beego.Router("/manager/category", &controllers.ManagerController{}, "post,get:Category")
	beego.Router("/manager/update-cate", &controllers.ManagerController{}, "get:UpdateCate")
	beego.Router("/manager/del-cate", &controllers.ManagerController{}, "get:DelCate")
	beego.Router("/manager/icon-cate", &controllers.ManagerController{}, "post:UpdateCateIcon")

	//login
	beego.Router("/login", &controllers.AccountController{}, "*:Login")
	beego.Router("/regist", &controllers.AccountController{}, "*:Regist")
	beego.Router("/doregist", &controllers.AccountController{}, "post:DoRegist")
	beego.Router("/logout", &controllers.AccountController{}, "*:Logout")

    //个人中心
	beego.Router("/user/:username", &controllers.UserController{}, "get:Index")                 //主页
	beego.Router("/user/:username/collection", &controllers.UserController{}, "get:Collection") //收藏
	beego.Router("/user/:username/follow", &controllers.UserController{}, "get:Follow")         //关注
	beego.Router("/user/:username/fans", &controllers.UserController{}, "get:Fans")             //粉丝
	beego.Router("/follow/:uid", &controllers.BaseController{}, "get:SetFollow")                //关注或取消关注
	beego.Router("/book/score/:id", &controllers.BookController{}, "*:Score")                   //评分
	beego.Router("/book/comment/:id", &controllers.BookController{}, "post:Comment")            //评论
}
