package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"xiyu.com/common"
	"xiyu.com/facada"
	"xiyu.com/infra"
)

func wxlogin(c *gin.Context) {
	wxlog := facada.WxLogin{}
	if err := c.BindJSON(&wxlog); err != nil {
		common.GlbBaInfa.Logger.Infoln("Can not find args")
		c.String(http.StatusBadRequest, "Can not find args")
		return
	}

	rsp, err := infra.Code2Session(wxlog.Code)
	if err != nil {
		common.Logger.Infof("Code2Session failed:%s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	//jwt
	jwt, err := common.GetToken(wxlog.Nickname, rsp.Openid)
	if err != nil {
		common.Logger.Warnf("jwt failed:%s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Authorization", fmt.Sprintf("Bearer %s", jwt))
	c.String(http.StatusOK, "ok")
}

func wxregister(c *gin.Context) {
	wxlog := facada.WxLogin{}
	if err := c.BindJSON(&wxlog); err != nil {
		common.GlbBaInfa.Logger.Infoln("Can not find args")
		c.String(http.StatusBadRequest, "Can not find args")
		return
	}

	rsp, err := infra.Code2Session(wxlog.Code)
	if err != nil {
		common.Logger.Infof("Code2Session failed:%s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	//jwt
	jwt, err := common.GetToken(wxlog.Nickname, rsp.Openid)
	if err != nil {
		common.Logger.Warnf("jwt failed:%s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Header("Authorization", fmt.Sprintf("Bearer %s", jwt))
	c.String(http.StatusOK, "ok")
}

func uploadAvator(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Find none file")
		return
	}
	common.Logger.Infof("filename:%s", file.Filename)
	openId := c.MustGet("openId")
	if openId == "" {
		c.String(http.StatusInternalServerError, "Internal Error")
		return
	}

	dst := fmt.Sprintf("%s/%s", common.GlbBaInfa.Conf.Http.Avator, openId)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Error")
		common.Logger.Infof("save upload file failed:%s", err.Error())
		return
	}

}

func getAvator(c *gin.Context) {
	openId := c.MustGet("openId")
	if openId == "" {
		c.String(http.StatusInternalServerError, "Internal Error")
		return
	}

	dst := fmt.Sprintf("%s/%s", common.GlbBaInfa.Conf.Http.Avator, openId)
	data, err := ioutil.ReadFile(dst)
	if err != nil {
		common.Logger.Infof("Read file error:%s", err.Error())
		c.String(http.StatusInternalServerError, "Internal Error")
		return
	}

	c.Data(http.StatusOK, "image/gif;image/jpeg;image/png", data)

}

func RegisterUserRest(router *gin.Engine, handler gin.HandlerFunc) {
	router.POST(fmt.Sprintf("%s/user/wxlogin", common.GlbBaInfa.Conf.Http.Prefix), wxlogin)
	router.POST(fmt.Sprintf("%s/user/wxregister", common.GlbBaInfa.Conf.Http.Prefix), wxregister)
	router.POST(fmt.Sprintf("%s/user/upload-avator", common.GlbBaInfa.Conf.Http.Prefix), handler, uploadAvator)
	router.GET(fmt.Sprintf("%s/user/get-avator", common.GlbBaInfa.Conf.Http.Prefix), handler, getAvator)
}
