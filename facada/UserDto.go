package facada

type User struct {
	Uid      string `json:"uid" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
}

type WxLogin struct {
	TID      string `json:"tid" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
}

type UserResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Usr  User   `json:"user"`
}
