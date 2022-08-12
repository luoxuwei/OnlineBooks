package initialize

import (
	_ "OnlineBooks/models" //调用models模块的init函数
)

func init() {
	sysinit()
    dbinit()
}