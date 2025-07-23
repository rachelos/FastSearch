package admin

import (
	"net/http"
	"net/url"
	"os"

	"gitee.com/rachel_os/fastsearch/web/admin/assets"
	"github.com/gin-gonic/gin"
)

func adminIndex(ctx *gin.Context) {
	file, err := assets.Static.ReadFile("dist/index.html")
	if err != nil && os.IsNotExist(err) {
		ctx.String(http.StatusNotFound, "not found")
		return
	}
	ctx.Data(http.StatusOK, "text/html", file)
}

func handlerStatic(c *gin.Context) {
	staticServer := http.FileServer(http.FS(assets.Static))
	c.Request.URL = &url.URL{Path: "dist" + c.Request.RequestURI}
	staticServer.ServeHTTP(c.Writer, c.Request)
}

func Register(router *gin.Engine, handlers ...gin.HandlerFunc) {
	//注册路由
	r := router.Group("/admin", handlers...)
	r.GET("/", adminIndex)
	router.GET("/assets/*filepath", handlerStatic)
}
