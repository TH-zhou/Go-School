package models

import (
	"github.com/astaxie/beego/orm"
)

type CourseClickCollect struct {
	Id           int `orm:"column(id);auto"`
	Corpid       int `orm:"column(corpid)" description:"商户id"`
	CourseId     int `orm:"column(course_id)" description:"课程id"`
	Clickcount   int `orm:"column(clickcount);null" description:"点击量"`
	Collectcount int `orm:"column(collectcount);null" description:"收藏量"`
}

func (t *CourseClickCollect) TableName() string {
	return "course_click_collect"
}

func init() {
	orm.RegisterModel(new(CourseClickCollect))
}
