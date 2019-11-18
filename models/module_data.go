package models

import (
	"github.com/astaxie/beego/orm"
)

type ModuleData struct {
	Id         int    `orm:"column(id);auto"`
	Corpid     int    `orm:"column(corpid)" description:"商户id"`
	ModuleId   int    `orm:"column(module_id);null" description:"模块id"`
	Name       string `orm:"column(name);size(10);null" description:"分类名"`
	Icon       string `orm:"column(icon);size(100);null" description:"图标"`
	DataId     int    `orm:"column(data_id);null" description:"关联数据id"`
	Sort       int    `orm:"column(sort);null" description:"排序"`
	Createtime int    `orm:"column(createtime);null" description:"创建时间"`
}

func (t *ModuleData) TableName() string {
	return "module_data"
}

func init() {
	orm.RegisterModel(new(ModuleData))
}
