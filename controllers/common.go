package controllers

import (
	"Go-School/models"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

// 接收的参数
var receiveMap map[string]interface{}

var AttributesStruct Attributes

// 可访问的属性
type Attributes struct {
	ParamMap map[string]interface{} // 接收的请求参数
	Userid   int                    // 用户userid
	Corpid   int                    // 商户id
}

// 查询总数所需结构体
type Count struct {
	Count int64
}

// 公共类
type CommonController struct {
	beego.Controller

	Attributes
}

func init() {
	// 目前用户id和商户id写死
	AttributesStruct.Userid = 5
	AttributesStruct.Corpid = 10

	// 注册mysql驱动相关
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:root@tcp(127.0.0.1:3306)/goschool?charset=utf8")
}

// 获取接收的参数
func (this *CommonController) getReceiveParam() {
	// json解析参数
	if unmarshalErr := json.Unmarshal(this.Ctx.Input.RequestBody, &receiveMap); unmarshalErr != nil {
		this.ReturnData(-1, map[string]interface{}{})
	} else {
		// 判断参数是否存在
		if paramMap, ok := receiveMap["request"].(map[string]interface{}); !ok {
			this.ReturnData(-1, map[string]interface{}{})
		} else {
			AttributesStruct.ParamMap = paramMap
		}
	}
}

// 返回数据
func (this *CommonController) ReturnData(code int, data map[string]interface{}) {
	returnMap := make(map[string]interface{})
	returnMap["errcode"] = code
	errmsgModel := models.Errmsg{}
	returnMap["errmsg"] = errmsgModel.GetErrmsgByCode(code)
	if data, ok := data["data"]; ok {
		returnMap["data"] = data
	} else {
		returnMap["data"] = map[string]interface{}{}
	}

	if count, ok := data["count"]; ok {
		returnMap["count"] = count
	} else {
		returnMap["count"] = 0
	}

	switch beego.AppConfig.String("other::returnType") {
	case "json":
		this.Data["json"] = returnMap
		this.ServeJSON()
	case "xml":
		this.Data["xml"] = returnMap
		this.ServeXML()
	case "jsonp":
		this.Data["jsonp"] = returnMap
		this.ServeJSONP()
	}
}

// 检查必传参数
func (this *CommonController) checkParams(params []string) map[string]interface{} {
	// 获取参数
	this.getReceiveParam()

	paramMap := AttributesStruct.ParamMap
	for _, param := range params {
		if val, ok := paramMap[param]; !ok {
			beego.Info(param)
			this.ReturnData(51002, map[string]interface{}{})
		} else {
			// 断言
			switch val.(type) {
			case string, int:
				var valString string
				if valInt, ok := val.(int); ok {
					valString = strconv.Itoa(valInt)
				} else {
					valString = val.(string)
				}
				if len(valString) == 0 {
					this.ReturnData(51002, map[string]interface{}{})
				}
			case interface{}:
				if val.(interface{}) == nil {
					this.ReturnData(51002, map[string]interface{}{})
				}
			case []interface{}:
				valInterSli := val.([]interface{})
				if len(valInterSli) == 0 {
					this.ReturnData(51002, map[string]interface{}{})
				}
			}
		}
	}

	return paramMap
}

// 格式化时间
func (this *CommonController) FormatTime(datetime int, level string) string {
	//
	var timeLayout string
	switch strings.ToLower(level) {
	case "s": // 精确到秒
		timeLayout = "2006-01-02 15:04:05"
	case "d": // 精确到天
		timeLayout = "2006-01-02"
	case "m": // 精确到分
		timeLayout = "2006-01-02 15:04"
	default:
		timeLayout = "2006-01-02 15:04:05"
	}
	//设置本地时区
	//local, _ := time.LoadLocation("Local")

	timeInt64 := int64(datetime)
	formatTimeString := time.Unix(timeInt64, 0).Format(timeLayout)

	return formatTimeString
}

// 将格式化好的时间转为时间戳
func (this *CommonController) FormtimeTimeByString(date string) int64 {
	// 转换模板
	timeLayout := "2006-01-02 15:04:05"
	// 获取时区
	loc, _ := time.LoadLocation("Local")
	// 使用模板在对应时区转化为time.time类型
	theTime, _ := time.ParseInLocation(timeLayout, date, loc)
	// 转换为时间戳 类型是int64
	sr := theTime.Unix()

	return sr
}

// 分页信息
func (this *CommonController) ReturnPageConf() (pageNum, pageCount int) {
	// 当前第N页
	if page_num, ok := AttributesStruct.ParamMap["page_num"]; !ok {
		pageNum = 1
	} else {
		pageNumInt, _ := strconv.Atoi(page_num.(string))
		pageNum = pageNumInt
	}

	// 显示数量
	if page_count, ok := AttributesStruct.ParamMap["page_count"]; !ok {
		pageCount = 10
	} else {
		pageCountInt, _ := strconv.Atoi(page_count.(string))
		pageCount = pageCountInt
	}

	pageNum = (pageNum - 1) * pageCount

	return
}

// 将sql处理成查询总数sql
func (this *CommonController) ReturnCountSql(sql string) string {
	query := strings.ToLower(sql)

	betweenString := this.GetBetweenStr(query, "from", "order by")

	countSql := "select count(*) count " + betweenString

	return countSql
}

// 查找字符串中间的内容
func (this *CommonController) GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}
