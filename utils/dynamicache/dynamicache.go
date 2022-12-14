package dynamicache

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"strconv"
	"time"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/gomodule/redigo/redis"
)

var (
	pool      *redis.Pool = nil
	MaxIdle   int         = 0
	MaxOpen   int         = 0
	ExpireSec int64       = 0
)

func InitCache() {
	addr, _ := beego.AppConfig.String("dynamicache_addrstr")
	if len(addr) == 0 {
		addr = "127.0.0.1:6379"
	}
	if MaxIdle <= 0 {
		MaxIdle = 256
	}
	pass, _ := beego.AppConfig.String("dynamicache_passwd")
	if len(pass) == 0 {
		pool = &redis.Pool{
			MaxIdle:     MaxIdle,
			MaxActive:   MaxOpen,
			IdleTimeout: time.Duration(120),
			Dial: func() (redis.Conn, error) {
				return redis.Dial(
					"tcp",
					addr,
					redis.DialReadTimeout(1*time.Second),
					redis.DialWriteTimeout(1*time.Second),
					redis.DialConnectTimeout(1*time.Second),
				)
			},
		}
	} else {
		pool = &redis.Pool{
			MaxIdle:     MaxIdle,
			MaxActive:   MaxOpen,
			IdleTimeout: time.Duration(120),
			Dial: func() (redis.Conn, error) {
				return redis.Dial(
					"tcp",
					addr,
					redis.DialReadTimeout(1*time.Second),
					redis.DialWriteTimeout(1*time.Second),
					redis.DialConnectTimeout(1*time.Second),
					redis.DialPassword(pass),
				)
			},
		}
	}
}

func rdsdo(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := pool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	parmas := make([]interface{}, 0)
	parmas = append(parmas, key)

	if len(args) > 0 {
		for _, v := range args {
			parmas = append(parmas, v)
		}
	}
	return con.Do(cmd, parmas...)
}

func WriteString(key string, value string) error {
	_, err := rdsdo("SET", key, value)
	logs.Debug("redis set:" + key + "-" + value)
	rdsdo("EXPIRE", key, ExpireSec)
	return err
}

func ReadString(key string) (string, error) {
	result, err := rdsdo("GET", key)
	logs.Debug("redis get:" + key)
	if nil == err {
		str, _ := redis.String(result, err)
		return str, nil
	} else {
		logs.Debug("redis get error:" + err.Error())
		return "", err
	}
}

func WriteStruct(key string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if nil == err {
		return WriteString(key, string(data))
	} else {
		return nil
	}
}
func ReadStruct(key string, obj interface{}) error {
	if data, err := ReadString(key); nil == err {
		return json.Unmarshal([]byte(data), obj)
	} else {
		return err
	}
}

func WriteList(key string, list interface{}, total int) error {
	realKeyList := key + "_list"
	realKeyCount := key + "_count"
	data, err := json.Marshal(list)
	if nil == err {
		WriteString(realKeyCount, strconv.Itoa(total))
		return WriteString(realKeyList, string(data))
	} else {
		return nil
	}
}

func ReadList(key string, list interface{}) (int, error) {
	realKeyList := key + "_list"
	realKeyCount := key + "_count"
	if data, err := ReadString(realKeyList); nil == err {
		totalStr, _ := ReadString(realKeyCount)
		total := 0
		if len(totalStr) > 0 {
			total, _ = strconv.Atoi(totalStr)
		}
		return total, json.Unmarshal([]byte(data), list)
	} else {
		return 0, err
	}
}
