package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	_ "iHome/models"
	_ "iHome/routers"
	"net/http"
	"strings"
)

func ignoreStaticPath() {

	//透明static
	beego.InsertFilter("/", beego.BeforeRouter, TransparentStatic)
	beego.InsertFilter("/*", beego.BeforeRouter, TransparentStatic)
}

func TransparentStatic(ctx *context.Context) {
	orpath := ctx.Request.URL.Path
	beego.Debug("request url: ", orpath)
	//如果请求uri还有api字段,说明是指令应该取消静态资源路径重定向
	if strings.Index(orpath, "api") >= 0 {
		return
	}
	http.ServeFile(ctx.ResponseWriter, ctx.Request, "static/html/"+ctx.Request.URL.Path)
}

func main() {

	// 设置fastdfs的请求的静态路径重定向
	beego.SetStaticPath("/group1/M00", "fastdfs/storage_data/data")
	/*
		// 测试上传一个文件到fastdfs中
		_, fileId, _ := models.FDFSUploadByFileName("./main.go")
		fmt.Println("==== fileid ====", fileId)
	*/
	ignoreStaticPath()
	beego.Run()
}
