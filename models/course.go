package models

import (
	"github.com/astaxie/beego/orm"
)

type Course struct {
	Id               int    `orm:"column(id);auto"`
	Corpid           int    `orm:"column(corpid)" description:"商户id"`
	CourseCategoryId int    `orm:"column(course_category_id);null" description:"课程类型id"`
	Name             string `orm:"column(name);size(30);null" description:"课程名称"`
	Img              string `orm:"column(img);size(100);null" description:"课程图片"`
	Intro            string `orm:"column(intro);size(100);null" description:"课程介绍"`
	IsNew            int8   `orm:"column(is_new);null" description:"标记为新上线 1：标记 0：不标记"`
	IsTop            int8   `orm:"column(is_top);null" description:"是否置顶 0：否 1：是"`
	Userid           int    `orm:"column(userid);null" description:"上传人"`
	Status           int8   `orm:"column(status);null" description:"状态 1：上架 0：下架"`
	IsAssign         int8   `orm:"column(is_assign);null" description:"0全部可见 1指定可见"`
	IsNologinView    int8   `orm:"column(is_nologin_view);null" description:"未登录可查看 1可查看 0不能查看"`
	Updatetime       int    `orm:"column(updatetime);null" description:"更新时间"`
	Createtime       int    `orm:"column(createtime);null" description:"学习截止时间"`
}

func (t *Course) TableName() string {
	return "course"
}

func init() {
	orm.RegisterModel(new(Course))
}
