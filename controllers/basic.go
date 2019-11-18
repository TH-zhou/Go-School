package controllers

import (
	"Go-School/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type BasicController struct {
	CommonController
}

// 热词搜索
func (this *BasicController) HotSearch() {
	var respsMap []map[string]string

	// 热词搜索
	hotsSlice := []string{"123", "456", "789"}
	for _, hot := range hotsSlice {
		respsMap = append(respsMap, map[string]string{"words": hot})
	}

	returnMap := make(map[string]interface{})
	returnMap["rows"] = respsMap

	dataMap := make(map[string]interface{})
	dataMap["data"] = returnMap

	this.CommonController.ReturnData(0, dataMap)
}

// 基础信息
func (this *BasicController) BaseInfo() {
	// 获取课程类型
	basicsModel := []models.Basics{}
	ormer := orm.NewOrm()
	cond := orm.NewCondition()

	if _, allErr := ormer.QueryTable(new(models.Basics)).SetCond(cond.And("type", "course_type").And("is_delete", 0).AndCond(cond.And("corpid", 0).Or("corpid", AttributesStruct.Corpid))).All(&basicsModel, "value", "alias"); allErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}
	var couseTypesMap []map[string]string
	for _, basic := range basicsModel {
		couseTypesMap = append(couseTypesMap, map[string]string{"id": basic.Value, "value": basic.Alias})
	}

	returnMap := map[string]interface{}{
		"data": map[string]interface{}{
			"ranking_url": "xxxx",
			"about_url":   "xxxx",
			"logo_url":    "xxxx",
			"course_type": couseTypesMap,
		},
	}

	this.ReturnData(0, returnMap)
}

// 更新检查
func (this *BasicController) Update() {
	// 检查必传参数
	paramMap := this.checkParams([]string{"app_type", "app_ver"})
	beego.Info(paramMap)

	returnMap := map[string]interface{}{
		"data": map[string]interface{}{
			"new_app_ver":    "1.0",
			"level":          "1",
			"download_url":   "xxx",
			"update_explain": "xxx",
		},
	}

	this.ReturnData(0, returnMap)
}

// App在线心跳
func (this *BasicController) Online() {
	usersModel := models.Users{}
	// 查询数据
	ormer := orm.NewOrm()
	err := ormer.QueryTable(new(models.Users)).Filter("status", 1).Filter("userid", AttributesStruct.Userid).Filter("corpid", AttributesStruct.Corpid).One(&usersModel, "id", "duration")
	if err != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}

	// 数据不存在
	if err == orm.ErrNoRows {
		this.ReturnData(57004, map[string]interface{}{})
	}

	// + 30
	usersModel.Duration += 30

	// update
	num, upErr := ormer.Update(&usersModel, "duration")
	if upErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}
	if num == 0 {
		this.ReturnData(57003, map[string]interface{}{})
	}

	this.ReturnData(0, map[string]interface{}{})
}
