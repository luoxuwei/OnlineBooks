package main

import (
	_ "OnlineBooks/routers"
	_ "OnlineBooks/initialize"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

