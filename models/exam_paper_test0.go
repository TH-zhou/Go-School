package models

import (
	"github.com/astaxie/beego/orm"
)

type ExamPaperTest struct {
	Id                 int    `orm:"column(id);auto"`
	Corpid             int    `orm:"column(corpid)" description:"商户id"`
	ExamPaperId        int    `orm:"column(exam_paper_id)" description:"试卷id"`
	QuestionId         int    `orm:"column(question_id)" description:"题库id"`
	Title              string `orm:"column(title);size(200);null" description:"题目标题"`
	TitleImg           string `orm:"column(title_img);size(100);null" description:"试题图片"`
	Type               int8   `orm:"column(type);null" description:"试题类型"`
	Answer             string `orm:"column(answer);size(10);null" description:"试题答案"`
	Analysis           string `orm:"column(analysis);size(150);null" description:"试题解析"`
	QuestionTestOption string `orm:"column(question_test_option);null" description:"试题选项"`
	Score              string `orm:"column(score);size(5);null" description:"试题分值"`
	Status             int8   `orm:"column(status);null" description:"状态 1正常 0删除"`
	Createtime         int    `orm:"column(createtime);null" description:"创建时间"`
}

func (t *ExamPaperTest) TableName() string {
	return "exam_paper_test0"
}

func init() {
	orm.RegisterModel(new(ExamPaperTest))
}
