package models

import (
	"github.com/astaxie/beego/orm"
)

type Rotation struct {
	Id         int    `orm:"column(id);auto"`
	Corpid     int    `orm:"column(corpid)" description:"商户id"`
	Title      string `orm:"column(title);size(20);null" description:"轮播图标题"`
	Img        string `orm:"column(img);size(100);null" description:"图片"`
	Type       int8   `orm:"column(type);null" description:"类型 0：无 1：课程类型 2：链接"`
	Content    string `orm:"column(content);null" description:"内容"`
	Sort       int    `orm:"column(sort);null" description:"排序"`
	Status     int8   `orm:"column(status);null" description:"1正常 0已删除"`
	Createtime int    `orm:"column(createtime);null" description:"创建时间"`
}

func (t *Rotation) TableName() string {
	return "rotation"
}

func init() {
	orm.RegisterModel(new(Rotation))
}