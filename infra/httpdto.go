package infra

type Code2SessionRsp struct {
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Openid     string `json:"openid"`
	Msg        string `json:"errmsg"`
	Code       int    `json:"errcode"`
}
