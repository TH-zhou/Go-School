package routers

import (
	"Go-School/controllers"

	"github.com/astaxie/beego"
)

func init() {
	// Basic
	beego.Router("/Basic/hotSearch", &controllers.BasicController{}, "post:HotSearch")
	beego.Router("/Basic/baseInfo", &controllers.BasicController{}, "post:BaseInfo")
	beego.Router("/Basic/update", &controllers.BasicController{}, "post:Update")
	beego.Router("/Basic/online", &controllers.BasicController{}, "get:Online")

	// 首页
	beego.Router("/HomePage/baseInfo", &controllers.HomePageController{}, "get:BaseInfo")

	// Course
	beego.Router("/Course/courseList", &controllers.CourseController{}, "get:CourseList")

	// Exam
	beego.Router("/Exam/examSubmit", &controllers.ExamController{}, "post:ExamSubmit")
}
