package infra

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin/binding"
	"xiyu.com/common"
)

type httpConnector struct {
	client *http.Client
}

var httpCon *httpConnector
var httpOnce sync.Once

func get() *httpConnector {
	httpOnce.Do(func() {
		tr := &http.Transport{
			MaxIdleConns:       10,
			MaxConnsPerHost:    100,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		}
		client := &http.Client{Transport: tr}
		httpCon = &httpConnector{client: client}
		common.TaddItem(httpCon)
	})
	return httpCon
}

func (hc *httpConnector) Close() {
}

func Code2Session(code string) (*Code2SessionRsp, error) {
	hc := get()
	urlTemp := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	url := fmt.Sprintf(urlTemp, common.GlbBaInfa.Conf.Xxc.AppId, common.GlbBaInfa.Conf.Xxc.Token, code)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		common.Logger.Warnf("NewRequest err:%s", err.Error())
		return nil, err
	}
	rsp, err := hc.client.Do(req)
	if err != nil {
		common.Logger.Warnf("Do Request err:%s", err.Error())
		return nil, err
	}
	defer rsp.Body.Close()
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		common.Logger.Warnf("Read Request err:%s", err.Error())
		return nil, err
	}

	csr := Code2SessionRsp{}
	err = binding.JSON.BindBody(body, &csr)
	if err != nil {
		common.Logger.Warnf("Parser Request err:%s", err.Error())
		return nil, err
	}

	return &csr, nil
}
