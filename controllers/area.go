package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"iHome/models"
	"time"
)

type AreaController struct {
	beego.Controller
}

func (this *AreaController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

// /api/v1.0/area [get]
func (this *AreaController) GetAreas() {
	beego.Info("=== area controller is called ===")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}

	defer this.RetData(&resp)

	// 1. 应该从缓存中取得地域信息数据 直接返回给前端
	cache_conn, err := cache.NewCache("redis", `{"key":"ihome_go", "conn":"127.0.0.1:6379", "dbNum":"0"}`)
	if err != nil {
		beego.Info("cache redis conn error, err =", err)
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(models.RECODE_DBERR)
		return
	}
	// 2) 尝试从缓存中取得areas数据
	areas_info_value := cache_conn.Get("area_info")
	if areas_info_value != nil {
		beego.Info("=== get area_info from cache !!! ===")

		// 将areas_info_value(字符串) --> area 结构体切片
		var area_info interface{}
		json.Unmarshal(areas_info_value.([]byte), &area_info)
		resp.Data = area_info
		return
	}
	// 3) 假如没有

	/*
		// 将键值对beegorediskey:hahaha 存入redis数据库 有效时间300s
		// 假如失败则返回失败
		if err := cache_conn.Put("beegorediskey", "hahaha", time.Second*300); err != nil {
			beego.Info("put cache error")
			resp.Errno = models.RECODE_DBERR
			resp.Errmsg = models.RecodeText(models.RECODE_DBERR)
			return
		}
		beego.Info("put redis succ!")
		// 获取redis数据库中键为beegorediskey的value
		value := cache_conn.Get("beegorediskey")
		// 打印value
		fmt.Printf("%s\n", value)
	*/

	// 2. 如果缓存没有数据, 从mysql中查询area数据
	// 创建orm句柄
	o := orm.NewOrm()

	var areas []models.Area

	qs := o.QueryTable("area")
	row, err := qs.All(&areas)
	if err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 3. 将area数据存入缓存
	// 将area数据变成json
	areas_info_str, _ := json.Marshal(areas)
	if err := cache_conn.Put("area_info", areas_info_str, time.Second*3600); err != nil {
		beego.Info("set area_info to cache error, err =", err)
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 4. 将area数据 变成json发送给前端
	if row == 0 {
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	beego.Info("areas =", areas)
	resp.Data = areas
	return
}
