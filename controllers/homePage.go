package controllers

import (
	"Go-School/models"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type HomePageController struct {
	CommonController
}

// 首页
func (this *HomePageController) BaseInfo() {
	returnMap := make(map[string]interface{})

	// 轮播图
	returnMap["banner"] = this.getRotationList()

	// 场景清洁
	returnMap["scene"] = this.getScene()

	// 热门活动
	returnMap["activities"] = this.getActivities()

	// 最新课程
	returnMap["new_course"] = this.getNewCourses()

	// 热门课程
	returnMap["hot_course"] = this.getHotCourse()

	dataMap := make(map[string]interface{})
	dataMap["data"] = returnMap

	this.ReturnData(0, dataMap)
}

// 轮播图
func (this *HomePageController) getRotationList() []map[string]string {
	defaultCorpidInt, _ := strconv.Atoi(beego.AppConfig.String("other::defaultCorpid"))
	corpidsSlice := []int{
		AttributesStruct.Corpid,
		defaultCorpidInt,
	}

	rotationMap := []map[string]string{}
	for _, corpid := range corpidsSlice {
		rotationMap = this.getCorpRotationList(corpid)
		if len(rotationMap) > 0 {
			break
		}
	}

	return rotationMap
}

// 需要轮播图数据
func (this *HomePageController) getCorpRotationList(corpid int) []map[string]string {
	var rotationsModel []models.Rotation

	var ormer = orm.NewOrm()
	if _, allErr := ormer.QueryTable(new(models.Rotation)).Filter("corpid", corpid).Filter("status", 1).
		OrderBy("-sort").Limit(6).All(&rotationsModel, "img", "content"); allErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}

	returnMap := []map[string]string{}
	for _, rotation := range rotationsModel {
		returnMap = append(returnMap, map[string]string{"image": rotation.Img, "url": rotation.Content})
	}

	return returnMap
}

func (this *HomePageController) getScene() []map[string]interface{} {
	defaultCorpidInt, _ := strconv.Atoi(beego.AppConfig.String("other::defaultCorpid"))
	corpidsSlice := []int{
		AttributesStruct.Corpid,
		defaultCorpidInt,
	}

	sceneMap := []map[string]interface{}{}
	for _, corpid := range corpidsSlice {
		sceneMap = this.findCorpModuleData(corpid, 2)
		if len(sceneMap) > 0 {
			break
		}
	}

	return sceneMap
}

// 获取场景清洁数据
func (this *HomePageController) findCorpModuleData(corpid, module_id int) []map[string]interface{} {
	ormer := orm.NewOrm()
	moduleDatasModel := []models.ModuleData{}
	if _, allErr := ormer.QueryTable(new(models.ModuleData)).Filter("corpid", corpid).Filter("module_id", module_id).All(&moduleDatasModel, "data_id", "name", "icon"); allErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}

	moduleDataMap := []map[string]interface{}{}
	for _, moduleData := range moduleDatasModel {
		moduleDataMap = append(moduleDataMap, map[string]interface{}{"param": moduleData.DataId, "title": moduleData.Name, "icon": moduleData.Icon})
	}

	return moduleDataMap
}

// 热门活动
func (this *HomePageController) getActivities() []map[string]string {
	ormer := orm.NewOrm()
	cond := orm.NewCondition()

	activitysModelSlice := []models.Activity{}
	if _, allErr := ormer.QueryTable(new(models.Activity)).Filter("status", 1).SetCond(cond.AndCond(cond.And("corpid", AttributesStruct.Corpid).Or("corpid", beego.AppConfig.String("other::defaultCorpid")))).
		OrderBy("-id").Limit(5).All(&activitysModelSlice, "name", "content"); allErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}

	activitysMap := []map[string]string{}
	for _, activity := range activitysModelSlice {
		activitysMap = append(activitysMap, map[string]string{"title": activity.Name, "url": activity.Content})
	}

	return activitysMap
}

type courseAndClickCollect struct {
	Id               int
	Corpid           int
	Name             string
	CourseCategoryId int
	Img              string
	Clickcount       int
}

// 获取最新课程
func (this *HomePageController) getNewCourses() []map[string]interface{} {
	ormer := orm.NewOrm()

	showCourseIdSlice := this.CorpShowCourseId()
	coursesModel := []models.Course{}
	if _, allErr := ormer.QueryTable(new(models.Course)).Filter("status", 1).Filter("is_assign", 0).Filter("id__in", showCourseIdSlice).
		OrderBy("-is_new", "-id").Limit(5).All(&coursesModel, "id", "name", "img", "updatetime"); allErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}
	coursesMap := []map[string]interface{}{}
	for _, course := range coursesModel {
		updatetime := this.FormatTime(course.Updatetime, "s")
		coursesMap = append(coursesMap, map[string]interface{}{"param": course.Id, "title": course.Name, "image": course.Img, "updatetime": updatetime})
	}

	return coursesMap
}

// 查询展示的课程Id
func (this *HomePageController) CorpShowCourseId() []int {
	ormer := orm.NewOrm()

	corpid := AttributesStruct.Corpid

	// 查找不展示的课程id
	corpNotShowCoursesMap := []orm.Params{}
	if _, corpNotShowCoursesValuesErr := ormer.QueryTable(new(models.CorpNotShowCourses)).Filter("corpid", corpid).Filter("status", 1).Values(&corpNotShowCoursesMap, "courseid"); corpNotShowCoursesValuesErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}

	var corpNotShowCoursesSlice []int
	for _, corpNotShowCourse := range corpNotShowCoursesMap {
		corpNotShowCoursesSlice = append(corpNotShowCoursesSlice, int(corpNotShowCourse["Courseid"].(int64)))
	}

	// 查询不展示的课程的其余的全部课程
	showCourseMap := []orm.Params{}
	placeHolder := this.getPlaceHolder(corpNotShowCoursesSlice)
	if _, showCourseErr := ormer.Raw("SELECT id FROM course WHERE id NOT IN("+placeHolder+") AND (corpid = ? OR corpid = ?)", corpNotShowCoursesSlice, AttributesStruct.Corpid, beego.AppConfig.String("other::defaultCorpid")).
		Values(&showCourseMap); showCourseErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	}

	var showCourseSlice []int
	for _, showCourse := range showCourseMap {
		courseidInt, _ := strconv.Atoi(showCourse["id"].(string))
		showCourseSlice = append(showCourseSlice, courseidInt)
	}

	return showCourseSlice
}

// 获取占位
func (this *HomePageController) getPlaceHolder(slic []int) (placeHolder string) {
	for i := 0; i < len(slic); i++ {
		if len(placeHolder) > 0 {
			placeHolder += ",?"
		} else {
			placeHolder = "?"
		}
	}

	return
}

type HotCourse struct {
	Id         int
	Img        string
	Name       string
	Updatetime int
	Clickcount int
}

// 获取热门课程
func (this *HomePageController) getHotCourse() []map[string]interface{} {
	ormer := orm.NewOrm()
	qb, _ := orm.NewQueryBuilder("mysql")

	// 需要展示的课程id
	showCourseId := this.CorpShowCourseId()
	// 获取占位
	placeHolder := this.getPlaceHolder(showCourseId)

	corpidString := strconv.Itoa(AttributesStruct.Corpid)

	qb.Select("a.id", "a.img", "a.name", "a.updatetime", "IFNULL(b.clickcount, 0) clicks").From("course a").
		LeftJoin("course_click_collect b").On("a.id = b.course_id AND b.corpid = ?").
		Where("a.id IN ("+placeHolder+") AND a.status = ? AND a.is_assign = ? AND (a.corpid = ? OR a.corpid = ?)").
		GroupBy("a.id").OrderBy("-b.clickcount", "-a.id").Limit(5)

	sqlString := qb.String()

	hotCourseStruct := []HotCourse{}
	if _, queryRowsErr := ormer.Raw(sqlString, corpidString, showCourseId, 1, 0, corpidString, beego.AppConfig.String("other::defaultCorpid")).QueryRows(&hotCourseStruct); queryRowsErr != nil {
		beego.Info(queryRowsErr)
		this.ReturnData(-1, map[string]interface{}{})
	}

	var courseMap []map[string]interface{}
	for _, hotCourse := range hotCourseStruct {
		updatetime := this.FormatTime(hotCourse.Updatetime, "s")
		courseMap = append(courseMap, map[string]interface{}{"param": hotCourse.Id, "image": hotCourse.Img, "title": hotCourse.Name, "updatetime": updatetime, "clicks": hotCourse.Clickcount, "state": "更新", "sta_color": "#C80000"})
	}

	return courseMap
}
