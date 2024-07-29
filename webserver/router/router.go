package router

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/edgehook/ithings/webserver/api"
	v1 "github.com/edgehook/ithings/webserver/api/v1"
	"github.com/edgehook/ithings/webserver/middlewares"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	if gin.Mode() == gin.DebugMode {
		r.Use(gin.Logger(), gin.Recovery())
	} else {
		r.Use(gin.Recovery())
	}

	r.Use(middlewares.Cors())
	r.POST("/login", api.Login)
	apiv1 := r.Group("/v1")
	//ai detect
	detect := apiv1.Group("aiDetect")
	detect.GET("/image/:filename", v1.GetDetectImage)
	detect.GET("/video/:filename", v1.GetDetectVideo)
	detect.POST("", v1.AddAiDetect)
	detect.DELETE("/:filename", v1.DeleteAiDetect)

	//support dashboard web
	var (
		vueAssetsRoutePath = "./frontend" // 前端编译出来的 dist 所在路径
	)
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	vueAssetsRoutePath = filepath.Join(dir, "frontend")
	r.StaticFile("/", path.Join(vueAssetsRoutePath, "index.html"))           // 指定资源文件 url.  127.0.0.1/ 这种
	r.StaticFile("/fav32.png", path.Join(vueAssetsRoutePath, "fav32.png"))   // 127.0.0.1/favicon.ico
	r.StaticFS("/assets", http.Dir(path.Join(vueAssetsRoutePath, "assets"))) // 以 assets 为前缀的 url
	r.StaticFS("/static", http.Dir(path.Join(vueAssetsRoutePath, "static")))
	return r
}
