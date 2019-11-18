package models

import (
	"github.com/astaxie/beego/orm"
)

type Activity struct {
	Id          int    `orm:"column(id);auto"`
	Corpid      int    `orm:"column(corpid)" description:"商户id"`
	Name        string `orm:"column(name);size(10);null" description:"活动名称"`
	Img         string `orm:"column(img);size(100);null" description:"活动封面"`
	Type        int8   `orm:"column(type);null" description:"活动类型 1自定义 2url"`
	Content     string `orm:"column(content);null" description:"活动内容"`
	Status      int8   `orm:"column(status);null" description:"状态 0：已过期 1：正常 2：已删除"`
	Createtime  int    `orm:"column(createtime);null" description:"创建时间"`
	Invalidtime int    `orm:"column(invalidtime);null" description:"失效时间"`
}

func (t *Activity) TableName() string {
	return "activity"
}

func init() {
	orm.RegisterModel(new(Activity))
}
