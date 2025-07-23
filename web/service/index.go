package service

import (
	"gitee.com/rachel_os/fastsearch/global"
	"gitee.com/rachel_os/fastsearch/searcher"
	"gitee.com/rachel_os/fastsearch/searcher/model"
)

type Index struct {
	Container *searcher.Container
}

func NewIndex() *Index {
	return &Index{
		Container: global.Container,
	}
}

// AddIndex 添加索引
func (i *Index) AddIndex(dbName string, request *model.IndexDoc) error {
	return i.Container.GetDataBase(dbName).IndexDocument(request)
}

// BatchAddIndex 批次添加索引
func (i *Index) BatchAddIndex(dbName string, documents []*model.IndexDoc) error {
	db := i.Container.GetDataBase(dbName)
	for _, doc := range documents {
		if err := db.IndexDocument(doc); err != nil {
			return err
		}
	}
	return nil
}

// RemoveIndex 删除索引
func (i *Index) RemoveIndex(dbName string, data *model.RemoveIndexModel) error {
	db := i.Container.GetDataBase(dbName)
	if err := db.RemoveIndex(data.Id); err != nil {
		return err
	}
	return nil
}
func (i *Index) IndexTaskCount(dbName string) int64 {
	db := i.Container.GetDataBase(dbName)
	return db.TempCount()
}
