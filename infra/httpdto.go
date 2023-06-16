package infra

type Code2SessionRsp struct {
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	Openid     string `json:"openid"`
	Msg        string `json:"errmsg"`
	Code       int    `json:"errcode"`
}

type grantGzhTokenRsp struct {
	Token   string `json:"access_token"`
	Expires int32  `json:"expires_in"`
	Msg     string `json:"errmsg"`
	Code    int    `json:"errcode"`
}

type NewsVo struct {
	Title   string `json:"title"`
	Tmid    string `json:"thumb_media_id"`
	Showpic int    `json:"show_cover_pic"`
	Author  string `json:"author"`
	Digest  string `json:"digest"`
	Content string `json:"content"`
	Url     string `json:"url"`
	CSUrl   string `json:"content_source_url"`
}

type ContentVo struct {
	NewsItem []NewsVo `json:"news_item"`
}

type MediaVo struct {
	Mid        string    `json:"media_id"`
	Cnt        ContentVo `json:"content"`
	UpdateTime int64     `json:"update_time"`
}

type BatchgetMaterialRsp struct {
	TotalCount int       `json:"total_count"`
	ItemCount  int       `json:"item_count"`
	Item       []MediaVo `json:"item"`
	Msg        string    `json:"errmsg"`
	Code       int       `json:"errcode"`
}
