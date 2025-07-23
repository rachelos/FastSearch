package controller

import (
	"gitee.com/rachel_os/fastsearch/web/service"
	"github.com/gin-gonic/gin"
)

// Notice 通知
func Notice(c *gin.Context) {
	noc := &service.Notice{}
	if err := c.ShouldBindJSON(&noc); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	if noc.Text == "" {
		ResponseErrorWithMsg(c, "请输入通知内容")
		return
	}
	r, _ := srv.Notice.Notice(service.Notice{
		Text:     noc.Text,
		ToUser:   noc.ToUser,
		ToUserID: noc.ToUserID,
	})
	ResponseSuccessWithData(c, r)
}
