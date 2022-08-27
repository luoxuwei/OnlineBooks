package controllers

import (
	"OnlineBooks/common"
	"OnlineBooks/models"
	"OnlineBooks/utils"
	"OnlineBooks/utils/pagecache"
	"compress/gzip"
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/session"
    _ "github.com/beego/beego/v2/server/web/session/redis"
	"golang.org/x/net/context"
	"io"
	"strings"
	"time"
)

type BaseController struct {
	beego.Controller
	Member          *models.Member    //用户
	Option          map[string]string //全局设置
	EnableAnonymous bool              //开启匿名访问
}

type CookieRemember struct {
	MemberId int
	Account  string
	Time     time.Time
}

var globalSessions *session.Manager

func init() {
	sessionConf, _ := beego.AppConfig.String("sessionproviderconfig")
	if len(sessionConf) == 0 {
		return
	}
	sessionConfig := &session.ManagerConfig{
		CookieName:"gosessionid",
		EnableSetCookie: true,
		Gclifetime:3600,
		Maxlifetime: 3600,
		Secure: false,
		CookieLifeTime: 3600,
		ProviderConfig: sessionConf,
	}
	globalSessions, _ = session.NewManager("redis",sessionConfig)
	go globalSessions.GC()
}

func (c *BaseController) Finish() {
	controllerName, actionName := c.GetControllerAndAction()

	if pagecache.NeedWrite(controllerName, actionName, c.Ctx.Input.Params()) {
		render, err := c.RenderString()
		if nil == err && len(render) > 0 {
			pagecache.Write(controllerName, actionName, &render, c.Ctx.Input.Params())
		}
	}
}

//每个子类Controller公用方法调用前，都执行一下Prepare方法
func (c *BaseController) Prepare() {

	//如果有缓存，则返回缓存内容
	controllerName, actionName := c.GetControllerAndAction()
	if pagecache.InCacheList(controllerName, actionName) {
		contentPtr, err := pagecache.Read(controllerName, actionName, c.Ctx.Input.Params())
		if nil == err && len(*contentPtr) > 0 {
			io.WriteString(c.Ctx.ResponseWriter, *contentPtr)
			logs.Debug(controllerName + "-" + actionName + "read cache")
			c.StopRun()
		}
	}

	c.Member = models.NewMember() //初始化
	c.EnableAnonymous = false
	//从session中获取用户信息
	if globalSessions == nil {
		c.Abort("500")
	}
	session, _ := globalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	defer session.SessionRelease(context.Background(), c.Ctx.ResponseWriter)
	memberInSession := false
	var member models.Member
	memberObj := session.Get(context.Background(), common.SessionName)
	if nil != memberObj {
		member = memberObj.(models.Member)
		if member.MemberId > 0 {
			memberInSession = true
		}
	}

	if memberInSession {
		c.Member = &member
	} else {
		//如果Cookie中存在登录信息，从cookie中获取用户信息
		if cookie, ok := c.GetSecureCookie(common.AppKey(), "login"); ok {
			var remember CookieRemember
			err := utils.Decode(cookie, &remember)
			if err == nil {
				member, err := models.NewMember().Find(remember.MemberId)
				if err == nil {
					c.SetMember(*member)
					c.Member = member
				}
			}
		}
	}
	if c.Member.RoleName == "" {
		c.Member.RoleName = common.Role(c.Member.MemberId)
	}
	c.Data["Member"] = c.Member
	c.Data["BaseUrl"] = c.BaseUrl()
	c.Data["SITE_NAME"] = "OnlineBooks"
	//设置全局配置
	c.Option = make(map[string]string)
	c.Option["ENABLED_CAPTCHA"] = "false"
}


// Ajax接口返回Json
func (c *BaseController) JsonResult(errCode int, errMsg string, data ...interface{}) {
	jsonData := make(map[string]interface{}, 3)
	jsonData["errcode"] = errCode
	jsonData["message"] = errMsg

	if len(data) > 0 && data[0] != nil {
		jsonData["data"] = data[0]
	}
	returnJSON, err := json.Marshal(jsonData)
	if err != nil {
		logs.Error(err)
	}
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	//启用gzip压缩
	if strings.Contains(strings.ToLower(c.Ctx.Request.Header.Get("Accept-Encoding")), "gzip") {
		c.Ctx.ResponseWriter.Header().Set("Content-Encoding", "gzip")
		w := gzip.NewWriter(c.Ctx.ResponseWriter)
		defer w.Close()
		w.Write(returnJSON)
		w.Flush()
	} else {
		io.WriteString(c.Ctx.ResponseWriter, string(returnJSON))
	}
	c.StopRun()
}

func (c *BaseController) BaseUrl() string {
	host, _ := beego.AppConfig.String("sitemap_host")
	if len(host) > 0 {
		if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
			return host
		}
		return c.Ctx.Input.Scheme() + "://" + host
	}
	return c.Ctx.Input.Scheme() + "://" + c.Ctx.Request.Host
}

// 设置登录用户信息
func (c *BaseController) SetMember(member models.Member) {
	session, _ := globalSessions.SessionStart(c.Ctx.ResponseWriter, c.Ctx.Request)
	defer session.SessionRelease(context.Background(), c.Ctx.ResponseWriter)
	if member.MemberId <= 0 {
		session.Delete(context.Background(), common.SessionName)
		session.Delete(context.Background(), "uid")
	} else {
		session.Set(context.Background(), common.SessionName, member)
		session.Set(context.Background(), "uid", member.MemberId)
	}
}

//关注或取消关注
func (c *BaseController) SetFollow() {
	if c.Member.MemberId == 0 {
		c.JsonResult(1, "请先登录")
	}
	uid, _ := c.GetInt(":uid")
	if uid == c.Member.MemberId {
		c.JsonResult(1, "不能关注自己")
	}
	cancel, _ := new(models.Fans).FollowOrCancel(uid, c.Member.MemberId)
	if cancel {
		c.JsonResult(0, "已成功取消关注")
	}
	c.JsonResult(0, "已成功关注")
}