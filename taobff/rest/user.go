package rest

import (
	"fmt"
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
	common.Logger.Infof("%+v", rsp)
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
	common.Logger.Infof("%+v", rsp)
	c.String(http.StatusOK, "ok")
}

func RegisterUserRest(router *gin.Engine) {
	router.POST(fmt.Sprintf("%s/user/wxlogin", common.GlbBaInfa.Conf.Http.Prefix), wxlogin)
	router.POST(fmt.Sprintf("%s/user/wxregister", common.GlbBaInfa.Conf.Http.Prefix), wxregister)
}
