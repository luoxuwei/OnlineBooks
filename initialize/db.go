package initialize


import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
)

//调用方式
//dbinit() 或 dbinit("w") 或 dbinit("default") //初始化主库
//dbinit("w","r")	//同时初始化主库和从库
//dbinit("w")
func dbinit(aliases ...string) {
	//如果是开发模式，则显示命令信息
	isDev := false
	if runmode,err := beego.AppConfig.String("runmode"); err != nil && runmode == "dev" {
		isDev = true
	}

	if len(aliases) > 0 {
		for _, alias := range aliases {
			registDatabase(alias)
			//主库 自动建表
			if "w" == alias {
				orm.RunSyncdb("default", false, isDev)
			}
		}
	} else {
		registDatabase("w")
		orm.RunSyncdb("default", false, isDev)
	}

	if isDev {
		orm.Debug = isDev
	}
}

func registDatabase(alias string) {
	if len(alias) == 0 {
		return
	}
	//连接名称
	dbAlias := alias
	if "w" == alias || "default" == alias {
		dbAlias = "default"
		alias = "w"
	}

	dbName, _:= beego.AppConfig.String("db_" + alias + "_database")
	dbUser, _ := beego.AppConfig.String("db_" + alias + "_username")
	dbPwd, _ := beego.AppConfig.String("db_" + alias + "_password")
	dbHost, _ := beego.AppConfig.String("db_" + alias + "_host")
	dbPort, _ := beego.AppConfig.String("db_" + alias + "_port")

	orm.RegisterDataBase(dbAlias, "mysql", dbUser+":"+dbPwd+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8mb4", 30)

}