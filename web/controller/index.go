package controller

import (
	"gitee.com/rachel_os/fastsearch/searcher/model"

	"github.com/gin-gonic/gin"
)

// AddIndex 添加索引
func AddIndex(c *gin.Context) {
	document := &model.IndexDoc{}
	if err := c.ShouldBindJSON(&document); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	dbName := c.Query("database")
	if dbName == "" {
		ResponseErrorWithMsg(c, "database is empty")
		return
	}
	err := srv.Index.AddIndex(dbName, document)
	if err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}

	ResponseSuccessWithData(c, nil)
}

// BatchAddIndex 批量添加索引
func BatchAddIndex(c *gin.Context) {
	documents := make([]*model.IndexDoc, 0)
	if err := c.BindJSON(&documents); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}

	dbName := c.Query("database")
	if dbName == "" {
		ResponseErrorWithMsg(c, "database is empty")
		return
	}

	err := srv.Index.BatchAddIndex(dbName, documents)
	if err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}

	ResponseSuccess(c)
}

// RemoveIndex 删除索引
func RemoveIndex(c *gin.Context) {
	removeIndexModel := &model.RemoveIndexModel{}
	if err := c.BindJSON(&removeIndexModel); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}

	dbName := c.Query("database")
	if dbName == "" {
		ResponseErrorWithMsg(c, "database is empty")
		return
	}

	if err := srv.Index.RemoveIndex(dbName, removeIndexModel); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}

	ResponseSuccess(c)
}

// IndexTask 索引任务
func IndexTaskCount(c *gin.Context) {
	dbName := c.Query("database")
	if dbName == "" {
		ResponseErrorWithMsg(c, "database is empty")
		return
	}
	ResponseSuccessWithData(c, srv.Index.IndexTaskCount(dbName))
}
