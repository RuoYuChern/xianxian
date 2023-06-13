package facada

type Article struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Desc   string `json:"desc"`
	Img    string `json:"img"`
	Url    string `json:"url"`
}

type BatchGetBlogRsp struct {
	Code  int        `json:"code"`
	Msg   string     `json:"msg"`
	Blogs []*Article `json:"blogs"`
}
