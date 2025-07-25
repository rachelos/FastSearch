package model

// Highlight 关键词高亮
type Highlight struct {
	PreTag  string `json:"preTag"`  //高亮前缀
	PostTag string `json:"postTag"` //高亮后缀
}
type Neg_Option struct {
	Query   bool `json:"query"`   //关键词拦截
	Content bool `json:"content"` //结果负面过滤
}

// SearchRequest 搜索请求
type SearchRequest struct {
	Query     string     `json:"query,omitempty" form:"database"`      // 搜索关键词
	Order     string     `json:"order,omitempty" form:"database"`      // 排序类型
	ScoreExp  string     `json:"scoreExp,omitempty" form:"scoreExp"`   // 分数计算表达式
	MaxLimit  int        `json:"maxlimit,omitempty" form:"maxlimit"`   // 最大结果集限制，默认1000
	FilterExp string     `json:"filterExp,omitempty" form:"filterExp"` // 过滤计算表达式
	Negative  Neg_Option `json:"negative,omitempty"`                   // 负面项过滤开关
	Page      int        `json:"page,omitempty" form:"database"`       // 页码
	Limit     int        `json:"limit,omitempty" form:"database"`      // 每页大小，最大1000，超过报错
	Highlight *Highlight `json:"highlight,omitempty" form:"database"`  // 关键词高亮
	Database  string     `json:"database" form:"database"`             // 数据库名字
}

func (s *SearchRequest) GetAndSetDefault() *SearchRequest {

	if s.Limit == 0 {
		s.Limit = 100
	}
	if s.Page == 0 {
		s.Page = 1
	}

	if s.Order == "" {
		s.Order = "desc"
	}

	return s
}

// SearchResult 搜索响应
type SearchResult struct {
	Time      float64       `json:"time,omitempty"`      //查询用时
	Total     int           `json:"total"`               //总数
	PageCount int           `json:"pageCount"`           //总页数
	Page      int           `json:"page,omitempty"`      //页码
	Limit     int           `json:"limit,omitempty"`     //页大小
	Documents []ResponseDoc `json:"documents,omitempty"` //文档
	Words     []string      `json:"words,omitempty"`     //搜索关键词
}
