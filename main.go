package main

import (
	_ "OnlineBooks/routers"
	_ "OnlineBooks/initialize"
	"OnlineBooks/utils/pagecache"
	"fmt"
	"github.com/beego/beego/v2/adapter/toolbox"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	task := toolbox.NewTask("clear_expired_cache", "*/2 * * * * *", func() error { fmt.Println("--delete cache---"); pagecache.ClearExpiredFiles(); return nil })
	toolbox.AddTask("mbook_task", task)
	toolbox.StartTask()
	defer toolbox.StopTask()
	beego.Run()
}

