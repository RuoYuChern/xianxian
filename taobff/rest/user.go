package rest

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"xiyu.com/common"
	"xiyu.com/facade"
	"xiyu.com/infra"
)

type httpGetWork struct {
	avator   string
	openid   string
	nickName string
	common.Action
}

func (h *httpGetWork) Call() error {
	// body, err := infra.GetWithUrl(h.avator)
	// if err != nil {
	// 	common.Logger.Infof("get user:%s, url:%s failed:%s", h.openid, h.avator, err.Error())
	// 	return err
	// }
	// fs, _ := infra.GetFs(infra.Aavator)
	// ref, err := fs.Write(&grpc.FsBlockVo{OriginId: h.openid, Content: body})
	// if err != nil {
	// 	common.Logger.Infof("write user:%s data failed:%s", h.openid, err.Error())
	// 	return err
	// }

	// usrDto := &grpc.UserDto{OpenId: h.openid, NickName: h.nickName, AvatorRef: ref}
	// err = infra.UserSave(usrDto)
	// if err != nil {
	// 	common.Logger.Infof("UserSave user:%s data failed:%s", h.openid, err.Error())
	// 	return err
	// }
	return nil
}

func wxlogin(c *gin.Context) {
	wxregister(c)
}

func wxregister(c *gin.Context) {
	wxlog := facade.WxLogin{}
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
	// 添加后台抓取 用户头像
	common.AddAaction(&httpGetWork{avator: wxlog.Avatar, openid: rsp.Openid, nickName: wxlog.Nickname})
}

func getAvator(c *gin.Context) {
	openId := c.MustGet("openId").(string)
	nickName := c.MustGet("Username").(string)

	if openId == "" || nickName == "" {
		c.String(http.StatusInternalServerError, "Internal Error")
		return
	}

	data, err := infra.UserGetAvator(openId)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Error")
		return
	}
	c.Data(http.StatusOK, "image/jpeg;image/png;image/gif", data)

}

func uploadAvator(c *gin.Context) {
	openId := c.MustGet("openId")
	nickName := c.MustGet("Username")

	common.GlbBaInfa.Logger.Infof("openId:%+v,nickName:%+v", openId, nickName)

	if openId == nil || nickName == nil {
		common.GlbBaInfa.Logger.Info("Can not find openId or nickName")
		c.String(http.StatusInternalServerError, "Find none user")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		common.GlbBaInfa.Logger.Info("Can not find file")
		c.String(http.StatusBadRequest, "Can not find file")
		return
	}
	src, err := file.Open()
	if err != nil {
		common.GlbBaInfa.Logger.Info("File open failed:%s", err.Error())
		c.String(http.StatusInternalServerError, "File error")
		return
	}
	defer src.Close()
	var b bytes.Buffer
	wr := bufio.NewWriter(&b)
	_, err = io.Copy(wr, src)
	if err != nil {
		common.GlbBaInfa.Logger.Info("Copy failed:%s", err.Error())
		c.String(http.StatusInternalServerError, "Copy error")
		return
	}

	err = infra.UserSaveAvator(nickName.(string), openId.(string), b.Bytes())
	if err != nil {
		common.GlbBaInfa.Logger.Infoln("Save failed:%s", err.Error())
		c.String(http.StatusInternalServerError, "Save error")
		return
	}
	c.String(http.StatusOK, "OK")
}

func RegisterUserRest(router *gin.Engine, handler gin.HandlerFunc) {
	router.POST(fmt.Sprintf("%s/user/wxlogin", common.GlbBaInfa.Conf.Http.Prefix), wxlogin)
	router.POST(fmt.Sprintf("%s/user/wxregister", common.GlbBaInfa.Conf.Http.Prefix), wxregister)
	router.GET(fmt.Sprintf("%s/user/get-avator", common.GlbBaInfa.Conf.Http.Prefix), handler, getAvator)
	router.POST(fmt.Sprintf("%s/user/upload-avator", common.GlbBaInfa.Conf.Http.Prefix), handler, uploadAvator)
}
