package models

import (
	"github.com/astaxie/beego/orm"
)

type ExamRecord struct {
	Id                 int    `orm:"column(id);auto"`
	Corpid             int    `orm:"column(corpid)" description:"商户id"`
	ExamId             int    `orm:"column(exam_id)" description:"考试id"`
	Userid             int    `orm:"column(userid)" description:"用户userid"`
	LastScore          string `orm:"column(last_score);size(5);null" description:"最后一次考试分数"`
	Status             int8   `orm:"column(status);null" description:"状态 0待考试 1待补考 2待阅卷 3未通过 4已通过"`
	AlreadyRepeatCount int8   `orm:"column(already_repeat_count);null" description:"已补考次数"`
	LastCreatetime     int    `orm:"column(last_createtime);null" description:"最后一次考试提交时间"`
}

func (t *ExamRecord) TableName() string {
	return "exam_record0"
}

func init() {
	orm.RegisterModel(new(ExamRecord))
}
