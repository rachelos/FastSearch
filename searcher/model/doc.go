package model

// IndexDoc 索引实体
type IndexDoc struct {
	Id           string                 `json:"id,omitempty"`           //索引ID
	Title        string                 `json:"title,omitempty"`        //标题
	Text         string                 `json:"text,omitempty"`         //索引内容
	Flag         string                 `json:"flag,omitempty"`         //标识
	Time         string                 `json:"time,omitempty"`         //时间
	Tags         string                 `json:"tags,omitempty"`         //标签
	Seo          map[string]interface{} `json:"seo,omitempty"`          //补充关键词SEO关键词
	Keys         map[string]interface{} `json:"keys,omitempty"`         //补充关键词支持 如：site:test.com查询
	Cut_Document bool                   `json:"cut_document,omitempty"` //是否对文档进行索引
	Document     map[string]interface{} `json:"document,omitempty"`     //索引文档
}

// StorageIndexDoc 文档对象
type StorageIndexDoc struct {
	*IndexDoc
	Keys []string `json:"keys,omitempty"`
}

type ResponseDoc struct {
	IndexDoc
	OriginalText  string   `json:"originalText,omitempty"`
	OriginalTitle string   `json:"originalTitle,omitempty"`
	Flag          string   `json:"flag,omitempty"`
	Time          string   `json:"time,omitempty"`
	Tags          string   `json:"tags,omitempty"`
	Score         int      `json:"score,omitempty"` //得分
	Keys          []string `json:"keys,omitempty"`
}

type RemoveIndexModel struct {
	Id string `json:"id,omitempty"`
}

type ResponseDocSort []ResponseDoc

func (r ResponseDocSort) Len() int {
	return len(r)
}

func (r ResponseDocSort) Less(i, j int) bool {
	return r[i].Score < r[j].Score
}

func (r ResponseDocSort) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
