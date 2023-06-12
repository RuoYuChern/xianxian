package facada

type User struct {
	Uid      string `json:"uid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type WxLogin struct {
	TID      string `json:"tid"`
	Code     string `json:"code"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type UserResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Usr  User   `json:"user"`
}
