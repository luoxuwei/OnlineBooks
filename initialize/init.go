package initialize

import (
	_ "OnlineBooks/models" //调用models模块的init函数
	"github.com/beego/beego/v2/client/orm"
)

func init() {
	sysinit()
	orm.RegisterDriver("mysql", orm.DRMySQL)
    dbinit()
}