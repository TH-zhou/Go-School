package models

import (
	"github.com/astaxie/beego/orm"
)

type ExamPaper struct {
	Id         int    `orm:"column(id);auto"`
	Corpid     int    `orm:"column(corpid)" description:"商户id"`
	Title      string `orm:"column(title);size(20);null" description:"试卷名"`
	Status     int8   `orm:"column(status);null" description:"状态 1正常 0删除"`
	IsReading  int8   `orm:"column(is_reading);null" description:"是否需要阅卷 1需要 0不需要"`
	Createtime int    `orm:"column(createtime);null" description:"创建时间"`
}

func (t *ExamPaper) TableName() string {
	return "exam_paper"
}

func init() {
	orm.RegisterModel(new(ExamPaper))
}
