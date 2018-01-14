package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"iHome/models"
	"path"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

// /api/v1.0/session [get]
/*
func (this *UserController) GetSession() {
	beego.Info("=== getsession controller is called ===")

	resp := Resp{Errno: models.RECODE_USERERR, Errmsg: models.RecodeText(models.RECODE_USERERR)}
	defer this.RetData(&resp)
	return
}
*/
// /api/v1.0/users [post]
func (this *UserController) Reg() {
	beego.Info("=== reg controller is called ===")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	// 1. 得到客户端传递的信息 json解析
	// request
	var regRequestMap = make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &regRequestMap)

	beego.Info("client ret request =", regRequestMap)

	// 2. 校验信息的合法性 (mobile password sms_code)
	if regRequestMap["mobile"] == "" || regRequestMap["password"] == "" || regRequestMap["sms_code"] == "" {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 3.入库user
	user := models.User{}
	user.Mobile = regRequestMap["mobile"].(string)
	user.Password_hash = regRequestMap["password"].(string)
	user.Name = regRequestMap["mobile"].(string)

	o := orm.NewOrm()
	id, err := o.Insert(&user)
	if err != nil {
		beego.Info("insert error =", err)
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	beego.Info("reg succ!!! user id =", id)

	// 4.将用户存储到session中
	this.SetSession("name", user.Name)
	this.SetSession("user_id", id)
	this.SetSession("mobile", user.Mobile)
	return
}

// /api/v1.0/session [get]
func (this *UserController) GetSessionName() {
	beego.Info("=== getSessionName controller is called ===")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	name := this.GetSession("name")
	if name == nil {
		resp.Errno = models.RECODE_USERERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	//nameData := Name{Name: name.(string)}
	//resp.Data = nameData
	nameMap := make(map[string]interface{})
	nameMap["name"] = name.(string)
	resp.Data = nameMap

	return
}

// /api/v1.0/sessions [post]
// 登陆
func (this *UserController) Login() {
	beego.Info("=== login controller is called ===")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	// 1. 得到请求数据
	// request
	var loginRequestMap = make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &loginRequestMap)

	beego.Info("client login request =", loginRequestMap)

	// 2. 校验数据合法性
	if loginRequestMap["mobile"] == "" || loginRequestMap["password"] == "" {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 3. 根据信息 查询数据mobile是否在表中
	var user models.User

	o := orm.NewOrm()
	qs := o.QueryTable("user")
	// 3.1 假如表中没有数据 错误
	if err := qs.Filter("mobile", loginRequestMap["mobile"].(string)).One(&user); err != nil {
		// 没有任何数据
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	// 3.2 假如有数据 则判断password是否相等
	// 比较密码
	if user.Password_hash != loginRequestMap["password"].(string) {
		// 密码不正确
		resp.Errno = models.RECODE_PWDERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	beego.Info("login succ!!! user id =", user.Id)

	// 4. 将用户存储到session中
	this.SetSession("name", user.Name)
	this.SetSession("user_id", user.Id)
	this.SetSession("mobile", user.Mobile)

	return
}

// /api/v1.0/sessios [delete]
func (this *UserController) DeleteSessionName() {
	beego.Info("=== deleteSessionName controller is called ===")
	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	this.DelSession("name")
	this.DelSession("user_id")
	this.DelSession("mobile")
	return
}

// /api/v1.0/user/avatar [post]
// 上传头像
func (this *UserController) UploadAvatar() {
	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	// 拿到用户的文件二进制数据
	file, header, err := this.GetFile("avatar")
	if err != nil {
		resp.Errno = models.RECODE_SERVERERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	// 创建一个file文件的buffer
	fileBuffer := make([]byte, header.Size)

	if _, err = file.Read(fileBuffer); err != nil {
		resp.Errno = models.RECODE_IOERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// user01.jpg
	suffix := path.Ext(header.Filename) // ---> ".jpg"
	groupName, fileId, err := models.FDFSUploadByBuffer(fileBuffer, suffix[1:])
	if err != nil {
		resp.Errno = models.RECODE_IOERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		beego.Info("upload file error, name =", header.Filename)
		return
	}

	beego.Info("fdfs upload file succ groupName =", groupName, "fileId =", fileId)

	// 通过session 得到当前用的user_id
	user_id := this.GetSession("user_id")

	user := models.User{Id: user_id.(int), Avatar_url: fileId}
	o := orm.NewOrm()

	if _, err := o.Update(&user, "avatar_url"); err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 将fileId拼接一个完整的url路径 + ip + port返回给前端
	avatar_url := "http://39.106.152.53/" + fileId

	url_map := make(map[string]interface{})
	url_map["avatar_url"] = avatar_url
	resp.Data = url_map

	return
}

// /api/v1.0/user [get]
// 获取用户用户名和手机号
func (this *UserController) GetUser() {
	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	user := models.User{}
	user.Name = this.GetSession("name").(string)
	user.Mobile = this.GetSession("mobile").(string)

	resp.Data = user

	return
}

// /api/v1.0/user/name [put]
// 设置用户名
func (this *UserController) SetUserName() {
	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	// 1. 首先得到请求中 被修改后的用户名
	var newName = make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &newName)

	beego.Info("client setusername request =", newName)

	// 2. 校验信息合法性
	if newName["name"].(string) == "" {
		resp.Errno = models.RECODE_REQERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 3. 更新数据
	o := orm.NewOrm()
	user_id := this.GetSession("user_id")
	user := models.User{Id: user_id.(int), Name: newName["name"].(string)}
	if _, err := o.Update(&user, "name"); err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	this.SetSession("name", newName["name"].(string))

	resp.Data = newName

	return
}
