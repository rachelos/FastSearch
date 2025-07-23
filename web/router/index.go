package router

import (
	"gitee.com/rachel_os/fastsearch/web/controller"

	"github.com/gin-gonic/gin"
)

// InitIndexRouter 索引路由
func InitIndexRouter(Router *gin.RouterGroup) {

	indexRouter := Router.Group("index")
	indexRouter.Use(controller.CheckExistsDB())
	{
		indexRouter.POST("", controller.AddIndex)               // 添加单条索引
		indexRouter.POST("batch", controller.BatchAddIndex)     // 批量添加索引
		indexRouter.POST("remove", controller.RemoveIndex)      // 删除索引
		indexRouter.Any("taskcount", controller.IndexTaskCount) // 索引任务
	}
}
