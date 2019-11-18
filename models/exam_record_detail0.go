package models

import (
	"github.com/astaxie/beego/orm"
)

type ExamRecordDetail struct {
	Id           int    `orm:"column(id);auto"`
	Corpid       int    `orm:"column(corpid)" description:"商户id"`
	ExamRecordId int    `orm:"column(exam_record_id)" description:"考试记录id"`
	Start        int    `orm:"column(start);null" description:"开始考试时间"`
	End          int    `orm:"column(end);null" description:"结束考试时间"`
	TimeCost     int    `orm:"column(time_cost);null" description:"答题用时"`
	Score        string `orm:"column(score);size(5);null" description:"考试分数"`
	IsRepeat     int8   `orm:"column(is_repeat);null" description:"是否是补考记录 0不是 1是"`
	Status       int8   `orm:"column(status);null" description:"当前考试状态 0未通过 1已通过 "`
}

func (t *ExamRecordDetail) TableName() string {
	return "exam_record_detail0"
}

func init() {
	orm.RegisterModel(new(ExamRecordDetail))
}
