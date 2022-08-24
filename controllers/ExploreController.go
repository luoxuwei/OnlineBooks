package controllers

import (
	"OnlineBooks/models"
	"OnlineBooks/utils"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"math"
	"strconv"
)

type ExploreController struct {
	BaseController
}

func (c *ExploreController) Index() {
	var (
		cid       int //分类id
		cate      models.Category
		urlPrefix = beego.URLFor("ExploreController.Index")
	)

	cidstr := c.Ctx.Input.Param(":cid")
	if len(cidstr) > 0 {
		if cid, _ = strconv.Atoi(cidstr); cid > 0 {
			cateModel := new(models.Category)
			cate = cateModel.Find(cid)
			c.Data["Cate"] = cate
		}
	}

	c.Data["Cid"] = cid
	c.TplName = "explore/index.html"

	pageIndex, _ := c.GetInt("page", 1)
	pageSize := 24

	books, totalCount, err := models.NewBook().HomeData(pageIndex, pageSize, cid)
	if err != nil {
		logs.Error(err)
		c.Abort("404")
	}

	if totalCount > 0 {
		urlSuffix := ""
		if cid > 0 {
			urlSuffix = urlSuffix + "&cid=" + strconv.Itoa(cid)
		}
		html := utils.NewPaginations(4, totalCount, pageSize, pageIndex, urlPrefix, urlSuffix)
		c.Data["PageHtml"] = html
	} else {
		c.Data["PageHtml"] = ""
	}

	c.Data["TotalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))
	c.Data["Lists"] = books
}
