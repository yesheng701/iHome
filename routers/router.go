package routers

import (
	"github.com/astaxie/beego"
	"iHome/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	// 处理显示地区的请求
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetAreas")
	// 对房屋首页展示的业务
	beego.Router("/api/v1.0/houses/index", &controllers.HousesIndexController{}, "get:HousesIndex")
	// 处理用户session的请求
	beego.Router("/api/v1.0/session", &controllers.UserController{}, "get:GetSessionName;delete:DeleteSessionName")
	// 处理用户注册的请求
	beego.Router("/api/v1.0/users", &controllers.UserController{}, "post:Reg")
	// 处理用户登陆的请求
	beego.Router("/api/v1.0/sessions", &controllers.UserController{}, "post:Login")
	beego.Router("/api/v1.0/user", &controllers.UserController{}, "get:GetUser")
	// 处理上传请求
	beego.Router("/api/v1.0/user/avatar", &controllers.UserController{}, "post:UploadAvatar")
	beego.Router("/api/v1.0/user/name", &controllers.UserController{}, "put:SetUserName")
}
