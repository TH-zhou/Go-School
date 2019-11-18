package models

import (
	"github.com/astaxie/beego/orm"
)

type CorpNotShowCourses struct {
	Id       int  `orm:"column(id);auto"`
	Corpid   int  `orm:"column(corpid)" description:"商户id"`
	Courseid int  `orm:"column(courseid)" description:"不需要展示的课程id"`
	Status   int8 `orm:"column(status)" description:"状态 1：正常 0：已删除"`
}

func (t *CorpNotShowCourses) TableName() string {
	return "corp_not_show_courses"
}

func init() {
	orm.RegisterModel(new(CorpNotShowCourses))
}
