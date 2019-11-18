package models

import (
	"github.com/astaxie/beego/orm"
)

type UserScore struct {
	Id         int    `orm:"column(id);auto"`
	Corpid     int    `orm:"column(corpid)" description:"商户id"`
	Userid     int    `orm:"column(userid)" description:"用户id"`
	Score      string `orm:"column(score);size(3);null" description:"积分值 前面有-表示减积分"`
	Type       int8   `orm:"column(type);null" description:"积分类型"`
	Info       string `orm:"column(info);size(50);null" description:"详细信息"`
	Day        int    `orm:"column(day);null" description:"时间 精确到天"`
	Createtime int    `orm:"column(createtime);null" description:"创建时间"`
}

func (t *UserScore) TableName() string {
	return "user_score0"
}

func init() {
	orm.RegisterModel(new(UserScore))
}
