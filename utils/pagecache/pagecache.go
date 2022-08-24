package pagecache

import (
	"context"
	"errors"
	"github.com/beego/beego/v2/core/logs"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/client/cache"
)

var (
	BasePath  string              = ""
	ExpireSec int64               = 0
	store     *cache.FileCache    = nil
	cacheMap  map[string]bool     = nil
)

func InitCache() {
	store = &cache.FileCache{CachePath: BasePath}
	pagecacheList, _ := beego.AppConfig.Strings("pagecache_list")

	//初始化静态化配置列表
	cacheMap = make(map[string]bool)
	for _, v := range pagecacheList {
		cacheMap[strings.ToLower(v)] = true
	}
}

func InCacheList(controllerName, actionName string) bool {
	keyname := cacheKey(controllerName, actionName)
	if f := cacheMap[keyname]; f {
		return f
	}
	return false
}

func NeedWrite(controllerName, actionName string) bool {
	if InCacheList(controllerName, actionName) {
		keyname := cacheKey(controllerName, actionName)
		content, err := store.Get(context.Background(), keyname)
		if err == nil {
			if v := content.(string); len(v) > 0 {
				return false
			} else {
				logs.Debug("need write :" + keyname)
				return true
			}
		}

	}
	return false
}

func Write(controllerName, actionName string, content *string) error {
	keyname := cacheKey(controllerName, actionName)
	if len(keyname) == 0 {
		return errors.New("未找到缓存key")
	}

	err := store.Put(context.Background(), keyname, *content, time.Duration(ExpireSec)*time.Second)

	return err
}

func Read(controllerName, actionName string) (*string, error) {
	keyname := cacheKey(controllerName, actionName)
	if len(keyname) == 0 {
		return nil, errors.New("未找到缓存key")
	}
	content, err := store.Get(context.Background(), keyname)
	if err == nil {
		v := content.(string)
		return &v, nil
	}

	return nil, err
}

func cacheKey(controllerName, actionName string) string {
	if len(controllerName) > 0 && len(actionName) > 0 {
		rtnstr := strings.ToLower(controllerName + "_" + actionName)
		return rtnstr
	}

	return ""
}
