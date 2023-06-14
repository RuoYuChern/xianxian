package infra

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin/binding"
	"xiyu.com/common"
)

type wxToken struct {
	token   string
	expires int32
	start   time.Time
}

type httpConnector struct {
	client   *http.Client
	gzhToken *wxToken
	mu       sync.Mutex
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
		httpCon = &httpConnector{client: client, mu: sync.Mutex{}}
		common.TaddItem(httpCon)
	})
	return httpCon
}

func (hc *httpConnector) Close() {
}

func Code2Session(code string) (*Code2SessionRsp, error) {
	urlTemp := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	url := fmt.Sprintf(urlTemp, common.GlbBaInfa.Conf.Xxc.AppId, common.GlbBaInfa.Conf.Xxc.Token, code)
	body, err := doGet(url)
	if err != nil {
		return nil, err
	}
	csr := Code2SessionRsp{}
	err = binding.JSON.BindBody(body, &csr)
	if err != nil {
		common.Logger.Warnf("Parser Request err:%s", err.Error())
		return nil, err
	}

	if csr.Code != 0 {
		common.Logger.Infof("Code2Session failed:%s", csr.Msg)
		return nil, errors.New(csr.Msg)
	}
	return &csr, nil
}

func BatchgetMaterial(lx string, offset int, count int) (*BatchgetMaterialRsp, error) {
	hc := get()
	err := grantGzhToken()
	if err != nil {
		return nil, err
	}
	urlTemp := "https://api.weixin.qq.com/cgi-bin/material/batchget_material?access_token=%s"
	url := fmt.Sprintf(urlTemp, hc.gzhToken.token)

	var obj = map[string]any{"type": "news", "count": count, "offset": offset}
	body, err := doPost(url, obj)
	if err != nil {
		return nil, err
	}

	brsp := BatchgetMaterialRsp{}
	err = binding.JSON.BindBody(body, &brsp)
	if err != nil {
		common.Logger.Warnf("Parser Request err:%s", err.Error())
		return nil, err
	}

	return &brsp, nil
}

func GetMaterial(mid string) ([]byte, error) {
	hc := get()
	err := grantGzhToken()
	if err != nil {
		return nil, err
	}
	urlTemp := "https://api.weixin.qq.com/cgi-bin/material/get_material?access_token=%s"
	url := fmt.Sprintf(urlTemp, hc.gzhToken.token)
	var obj = map[string]string{"media_id": mid}
	body, err := doPost(url, obj)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func grantGzhToken() error {
	hc := get()
	for i := 0; i < 2; i++ {
		if hc.gzhToken != nil {
			gzhToken := hc.gzhToken
			n := time.Now()
			expire := gzhToken.start.Add(time.Duration(gzhToken.expires) * time.Second)
			if n.Before(expire) {
				return nil
			}
		}
		common.Logger.Info("Get gzh token")
		urlTemp := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
		url := fmt.Sprintf(urlTemp, common.GlbBaInfa.Conf.Gzh.AppId, common.GlbBaInfa.Conf.Gzh.Token)
		body, err := doGet(url)
		if err != nil {
			common.Logger.Warnf("Get gzh token error:%s", err.Error())
			return err
		}

		trsp := grantGzhTokenRsp{}
		err = binding.JSON.BindBody(body, &trsp)
		if err != nil {
			common.Logger.Warnf("Parser Request err:%s", err.Error())
			return err
		}
		if trsp.Code != 0 {
			common.Logger.Warnf("grant gzh token failed:%s", trsp.Msg)
			return err
		}
		common.Logger.Infof("token:%s", trsp.Token)
		if hc.gzhToken != nil {
			hc.gzhToken.start = time.Now()
			hc.gzhToken.token = trsp.Token
		} else {
			hc.gzhToken = &wxToken{token: trsp.Token, start: time.Now(), expires: trsp.Expires}
		}
	}
	return nil
}

func doPost(url string, obj any) ([]byte, error) {
	hc := get()
	js, err := json.Marshal(obj)
	if err != nil {
		common.Logger.Warnf("marshal failed:%s", err.Error())
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(js))
	if err != nil {
		common.Logger.Warnf("NewRequest:%s", err.Error())
		return nil, err
	}
	if err != nil {
		common.Logger.Warnf("Do Request err:%s", err.Error())
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
	return body, nil
}

func doGet(url string) ([]byte, error) {
	hc := get()
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
	return body, nil
}
