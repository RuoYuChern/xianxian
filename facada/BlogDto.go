package facada

type Article struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Img   string `json:"img"`
	Url   string `json:"url"`
}
