package models

import (
	"github.com/astaxie/beego/orm"
)

type Exam struct {
	Id             int    `orm:"column(id);auto"`
	Corpid         int    `orm:"column(corpid)" description:"商户id"`
	ExamPaperId    int    `orm:"column(exam_paper_id)" description:"试卷id"`
	CourseId       int    `orm:"column(course_id);null" description:"课程id"`
	CourseAssignId int    `orm:"column(course_assign_id);null" description:"课程学习任务id"`
	Name           string `orm:"column(name);size(30);null" description:"考试名称"`
	Img            string `orm:"column(img);size(100);null" description:"考试封面"`
	Type           int8   `orm:"column(type);null" description:"考试性质 1必考 2非必考"`
	Integral       int    `orm:"column(integral)" description:"考试通过积分"`
	Starttime      int    `orm:"column(starttime);null" description:"考试开始时间"`
	Endtime        int    `orm:"column(endtime);null" description:"考试结束时间"`
	Duration       int    `orm:"column(duration);null" description:"考试时长 单位s"`
	PassScore      int    `orm:"column(pass_score);null" description:"及格分数"`
	RepeatCount    int8   `orm:"column(repeat_count);null" description:"补考次数"`
	Status         int8   `orm:"column(status);null" description:"状态 1未开始 2进行中 3已结束 4已撤回"`
	ViewType       int8   `orm:"column(view_type);null" description:"查看答题详情 1完成后查看 2通过后查看"`
	CreateUserid   int    `orm:"column(create_userid);null" description:"创建考试的人"`
	Did            string `orm:"column(did);size(80);null" description:"分配的部门id，多个逗号隔开"`
	PostId         string `orm:"column(post_id);size(80);null" description:"分配的岗位id"`
	EntryStarttime int    `orm:"column(entry_starttime);null" description:"筛选入职开始时间"`
	EntryEndtime   int    `orm:"column(entry_endtime);null" description:"入职筛选结束时间"`
	Createtime     int    `orm:"column(createtime);null" description:"创建时间"`
}

func (t *Exam) TableName() string {
	return "exam"
}

func init() {
	orm.RegisterModel(new(Exam))
}
