package controllers

import (
	"OnlineBooks/common"
	"OnlineBooks/models"
	"OnlineBooks/utils"
	"OnlineBooks/utils/dynamicache"
	beego "github.com/beego/beego/v2/server/web"
	"strconv"
)

type UserController struct {
	BaseController
	UcenterMember models.Member //哪个用户的个人中心，你可以查看其他人的个人中心，所以不一定和登录用户一样
}

func (c *UserController) Prepare() {
	c.BaseController.Prepare()

	username := c.GetString(":username")

	// c.UcenterMember, _ = new(models.Member).GetByUsername(username)
	//读写缓存
	cachekeyUser := "dynamcache_user:" + username
	err := dynamicache.ReadStruct(cachekeyUser, &c.UcenterMember)
	if nil != err {
		c.UcenterMember, _ = new(models.Member).GetByUsername(username)
		dynamicache.WriteStruct(cachekeyUser, c.UcenterMember)
	}

	if c.UcenterMember.MemberId == 0 {
		c.Abort("404")
		return
	}
	c.Data["IsSelf"] = c.UcenterMember.MemberId == c.Member.MemberId
	c.Data["User"] = c.UcenterMember
	c.Data["Tab"] = "share"
}

//首页
func (c *UserController) Index() {
	page, _ := c.GetInt("page")
	pageSize := 10
	if page < 1 {
		page = 1
	}

	//从缓存读取c.Data["Books"]信息
	var books []*models.BookData
	cachekeyBookList := "dynamcache_userbook_" + strconv.Itoa(c.UcenterMember.MemberId) + "_page_" + strconv.Itoa(page)
	totalCount, err := dynamicache.ReadList(cachekeyBookList, &books)
	if nil != err {
		books, totalCount, _ = models.NewBook().SelectPage(page, pageSize, c.UcenterMember.MemberId, 0)
		dynamicache.WriteList(cachekeyBookList, books, totalCount)
	}
	c.Data["Books"] = books

	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, totalCount, pageSize, page, beego.URLFor("UserController.Index", ":username", c.UcenterMember.Account), "")
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Total"] = totalCount
	c.TplName = "user/index.html"
}

//收藏
func (c *UserController) Collection() {
	page, _ := c.GetInt("page")
	pageSize := 10
	if page < 1 {
		page = 1
	}

	//读取c.Data["Books"]信息
	var books []models.CollectionData
	var totalCount int64
	cachekeyCollectionList := "dynamcache_usercollection_" + strconv.Itoa(c.UcenterMember.MemberId) + "_page_" + strconv.Itoa(page)
	total, err := dynamicache.ReadList(cachekeyCollectionList, &books)
	totalCount = int64(total)
	if nil != err {
		totalCount, books, _ = new(models.Collection).List(c.UcenterMember.MemberId, page, pageSize)
		dynamicache.WriteList(cachekeyCollectionList, books, int(totalCount))
	}
	c.Data["Books"] = books

	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, int(totalCount), pageSize, page, beego.URLFor("UserController.Collection", ":username", c.UcenterMember.Account), "")
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Total"] = totalCount
	c.Data["Tab"] = "collection"
	c.TplName = "user/collection.html"
}

//关注
func (c *UserController) Follow() {
	page, _ := c.GetInt("page")
	pageSize := 18
	if page < 1 {
		page = 1
	}

	//读取关注列表缓存
	var fans []models.FansData
	var totalCount int64
	cachekeyfollowList := "dynamcache_userfollow_" + strconv.Itoa(c.UcenterMember.MemberId) + "_page_" + strconv.Itoa(page)
	total, err := dynamicache.ReadList(cachekeyfollowList, &fans)
	totalCount = int64(total)
	if nil != err { //数据库读取列表并缓存
		fans, totalCount, _ = new(models.Fans).FollowList(c.UcenterMember.MemberId, page, pageSize)
		dynamicache.WriteList(cachekeyfollowList, fans, int(totalCount))
	}

	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, int(totalCount), pageSize, page, beego.URLFor("UserController.Follow", ":username", c.UcenterMember.Account), "")
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Fans"] = fans
	c.Data["Tab"] = "follow"
	c.TplName = "user/fans.html"
}

//粉丝
func (c *UserController) Fans() {
	page, _ := c.GetInt("page")
	pageSize := 18
	if page < 1 {
		page = 1
	}

	// fans, totalCount, _ = new(models.Fans).FansList(c.UcenterMember.MemberId, page, pageSize)
	var fans []models.FansData
	var totalCount int64
	cachekeyFansList := "dynamcache_userfans_" + strconv.Itoa(c.UcenterMember.MemberId) + "_page_" + strconv.Itoa(page)
	total, err := dynamicache.ReadList(cachekeyFansList, &fans)
	totalCount = int64(total)
	if nil != err {
		fans, totalCount, _ = new(models.Fans).FansList(c.UcenterMember.MemberId, page, pageSize)
		dynamicache.WriteList(cachekeyFansList, fans, int(totalCount))
	}
	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, int(totalCount), pageSize, page, beego.URLFor("UserController.Fans", ":username", c.UcenterMember.Account), "")
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}
	c.Data["Fans"] = fans
	c.Data["Tab"] = "fans"
	c.TplName = "user/fans.html"
}