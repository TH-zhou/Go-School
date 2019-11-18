package controllers

import (
	"Go-School/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type CourseController struct {
	CommonController
}

type Course struct {
	Id               int    `json:"id"`
	Corpid           int    `json:"corpid"`
	Name             string `json:"name"`
	CourseType       string `json:"course_type"`
	CourseCategoryId int    `json:"course_category_id"`
	Img              string `json:"img"`
	Clickcount       int    `json:"clickcount"`
	Collection       bool   `json:"collection"`
}

func (this *CourseController) CourseList() {
	ormer := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")

	paramMap := AttributesStruct.ParamMap

	// 获取需要展示的课程Id
	var homePageController HomePageController
	showCourseSlice := homePageController.CorpShowCourseId()
	// 获取占位
	placeHolder := homePageController.getPlaceHolder(showCourseSlice)

	// 条件
	var whereString string
	whereString = "a.id IN (" + placeHolder + ") AND a.status = ? AND a.is_assign = ?"

	// 关键词搜索
	if keywords, ok := paramMap["keywords"]; ok {
		whereString += " AND a.name LIKE %" + keywords.(string) + "%"
	}

	// 课程类型搜索
	if course_type, ok := paramMap["course_type"]; ok {
		whereString += " AND a.course_category_id = " + course_type.(string)
	}

	// 默认企业和公共的都查询出来
	corpidString := strconv.Itoa(AttributesStruct.Corpid)
	defaultCorpidString := beego.AppConfig.String("other::defaultCorpid")
	defaultCorpidWhereString := " (a.corpid = " + corpidString + " OR a.corpid = " + defaultCorpidString + ")"
	if course_source, ok := paramMap["course_source"]; ok {
		courseSourceVal := course_source.(string)
		if courseSourceVal == "private" {
			defaultCorpidWhereString = " a.corpid = " + corpidString
		} else {
			defaultCorpidWhereString = " a.corpid = " + defaultCorpidString
		}
	}
	whereString += " AND " + defaultCorpidWhereString

	// 默认根据点击量排序
	orderByString := "b.clickcount"
	if sort_mode, ok := paramMap["sort_mode"]; ok && sort_mode.(string) == "time" {
		orderByString = "a.updatetime"
	}

	// 分页信息
	page_num, page_count := this.ReturnPageConf()

	qb.Select("a.id", "a.corpid", "a.name", "a.course_category_id", "a.img", "IFNULL(b.clickcount, 0) clickcount").
		From("course a").LeftJoin("course_click_collect b").On("a.id = b.course_id AND b.corpid = ?").
		Where(whereString).OrderBy(orderByString).Desc().Limit(page_count).Offset(page_num)

	sqlString := qb.String()

	var courses []Course
	if _, queryRowsErr := ormer.Raw(sqlString, corpidString, showCourseSlice, 1, 0).QueryRows(&courses); queryRowsErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}

	// 查询总数
	var count Count
	countSql := this.ReturnCountSql(sqlString)
	if queryRowErr := ormer.Raw(countSql, corpidString, showCourseSlice, 1, 0).QueryRow(&count); queryRowErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}

	var basicModel models.Basics
	for _, course := range courses {
		basic, basicNil := basicModel.GetBasicsTypeByvalue(course.Corpid, "course_type", course.CourseCategoryId)
		if basicNil != nil {
			this.ReturnData(-1, map[string]interface{}{})
		}
		course.CourseType = basic.Alias

		// 当前用户对该课程的收藏状态
		if collectCount, collectCountErr := ormer.QueryTable(new(models.CourseCollect)).Filter("corpid", AttributesStruct.Corpid).Filter("userid", AttributesStruct.Userid).Filter("course_id", course.Id).
			Filter("status", 1).Count(); collectCountErr != nil {
			this.ReturnData(-1, map[string]interface{}{})
		} else {
			var collectBool bool
			if collectCount > 0 {
				collectBool = true
			} else {
				collectBool = false
			}
			course.Collection = collectBool
		}
	}

	returnMap := map[string]interface{}{
		"data":  courses,
		"count": count.Count,
	}

	this.ReturnData(0, returnMap)
}
