package searcher

import (
	"fmt"
	"testing"

	"gitee.com/rachel_os/fastsearch/searcher/model"
)

func TestContainer_Init(t *testing.T) {
	c := &Container{
		Dir:   "./data/db",
		Debug: true,
	}
	err := c.Init()
	if err != nil {
		panic(err)
	}

	test := c.GetDataBase("default")

	fmt.Println(test.GetIndexCount())
	fmt.Println(test.GetDocumentCount())

	all := c.GetDataBases()
	for name, engine := range all {
		fmt.Println(name)
		fmt.Println(engine)
	}

	myMap := make(map[string]interface{})
	keys := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		// 向 map 中添加数据
		myMap["name"] = fmt.Sprintf("%s%v", "John Doe", i)
		myMap["age"] = 30 * i
		myMap["isStudent"] = false
		myMap["grades"] = []float64{90.5, 85.7, 92.3} // 可以存储任意类型，这里是一个 float64 切片
		keys["site"] = "www.google.com"
		test.IndexDocument(&model.IndexDoc{
			Id:           fmt.Sprintf("%s%v", "sn-", i),
			Text:         "hello world",
			Tags:         "abc",
			Document:     myMap,
			Seo:          keys,
			Keys:         keys,
			Cut_Document: true,
			Title:        "hello world",
		})
	}

}
