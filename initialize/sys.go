package initialize

import (
	"OnlineBooks/common"
	beego "github.com/beego/beego/v2/server/web"
	"path/filepath"
	"strings"
)

func sysinit() {
	//uploads静态路径
	uploads := filepath.Join(common.WorkingDirectory, "uploads")
	//localhost:8080/uploads
	beego.BConfig.WebConfig.StaticDir["/uploads"] = uploads

	registerFunctions()
}

//注册view里调用的函数
func registerFunctions() {
	beego.AddFuncMap("cdnjs", func(p string) string {
		cdn := beego.AppConfig.DefaultString("cdnjs", "")
		if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
			return cdn + string(p[1:])
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
			return cdn + "/" + p
		}
		return cdn + p
	})
	beego.AddFuncMap("cdncss", func(p string) string {
		cdn := beego.AppConfig.DefaultString("cdncss", "")
		if strings.HasPrefix(p, "/") && strings.HasSuffix(cdn, "/") {
			return cdn + string(p[1:])
		}
		if !strings.HasPrefix(p, "/") && !strings.HasSuffix(cdn, "/") {
			return cdn + "/" + p
		}
		return cdn + p
	})
}