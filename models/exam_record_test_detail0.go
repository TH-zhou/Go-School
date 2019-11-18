package models

import (
	"github.com/astaxie/beego/orm"
)

type ExamRecordTestDetail struct {
	Id                 int    `orm:"column(id);auto"`
	Corpid             int    `orm:"column(corpid)" description:"商户id"`
	ExamRecordDetailId int    `orm:"column(exam_record_detail_id)" description:"考试记录详情id"`
	ExamPaperTestId    int    `orm:"column(exam_paper_test_id)" description:"试卷下题目id"`
	ExamPaperTestType  int8   `orm:"column(exam_paper_test_type);null" description:"试卷下试题类型"`
	Answer             string `orm:"column(answer);null" description:"用户答案"`
	IsSuccess          int8   `orm:"column(is_success);null" description:"正确与否 0：错误 1：正确 2：无需知道对错（eg：问答题）"`
	Score              string `orm:"column(score);size(5);null" description:"用户得分"`
}

func (t *ExamRecordTestDetail) TableName() string {
	return "exam_record_test_detail0"
}

func init() {
	orm.RegisterModel(new(ExamRecordTestDetail))
}
