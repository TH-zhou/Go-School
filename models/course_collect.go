package models

import (
	"github.com/astaxie/beego/orm"
)

type CourseCollect struct {
	Id         int  `orm:"column(id);auto"`
	Corpid     int  `orm:"column(corpid)" description:"商户id"`
	Userid     int  `orm:"column(userid)" description:"用户id"`
	CourseId   int  `orm:"column(course_id);null" description:"课程id"`
	Status     int8 `orm:"column(status);null" description:"状态 1：已收藏 0：未收藏"`
	Createtime int  `orm:"column(createtime);null" description:"收藏时间"`
}

func (t *CourseCollect) TableName() string {
	return "course_collect0"
}

func init() {
	orm.RegisterModel(new(CourseCollect))
}
