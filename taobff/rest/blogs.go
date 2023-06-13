package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"xiyu.com/common"
	"xiyu.com/facada"
	"xiyu.com/infra"
)

func batchGet(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("page-size")

	if (page == "") || pageSize == "" {
		c.String(http.StatusBadRequest, "page or pageSize is error")
	}

	iPage, err := strconv.Atoi(page)
	if err != nil {
		c.String(http.StatusBadRequest, "page is error")
	}

	iPageSize, err := strconv.Atoi(pageSize)
	if err != nil {
		c.String(http.StatusBadRequest, "page is error")
	}
	rsp, err := infra.BatchgetMaterial("news", (iPage * iPageSize), iPageSize)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	blogs := make([]*facada.Article, iPageSize)
	for v := range rsp.Item {
		mvo := rsp.Item[v]
		news := mvo.Cnt.NewsItem[0]
		img := fmt.Sprintf("/bff/blogs/material-get?media_id=%s", news.Tmid)
		wz := facada.Article{Id: mvo.Mid, Type: "gzh", Title: news.Title,
			Author: news.Author, Desc: news.Digest, Img: img, Url: news.Url}
		blogs = append(blogs, &wz)
	}
	bgr := facada.BatchGetBlogRsp{Code: rsp.Code, Msg: rsp.Msg, Blogs: blogs}

	c.JSON(http.StatusOK, &bgr)
}

func getMaterial(c *gin.Context) {
	mid := c.Query("media_id")
	rsp, err := infra.GetMaterial(mid)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	data, err := common.ImgResize(rsp, 64, 64)
	if err != nil {
		c.Data(http.StatusOK, "image/gif;image/jpeg;image/png", rsp)
	} else {
		c.Data(http.StatusOK, "image/gif;image/jpeg;image/png", data)
	}

}

func RegisterBlogRest(router *gin.Engine) {
	router.GET(fmt.Sprintf("%s/blogs/batch-get", common.GlbBaInfa.Conf.Http.Prefix), batchGet)
	router.GET(fmt.Sprintf("%s/blogs/material-get", common.GlbBaInfa.Conf.Http.Prefix), getMaterial)
}
