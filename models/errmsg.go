package models

import (
	"github.com/astaxie/beego/orm"
)

type Errmsg struct {
	Id   int    `orm:"column(id);auto"`
	Code int    `orm:"column(code);null" description:"错误码"`
	Msg  string `orm:"column(msg);size(20);null" description:"错误详情"`
}

func (t *Errmsg) TableName() string {
	return "errmsg"
}

func init() {
	orm.RegisterModel(new(Errmsg))
}

// GetErrmsgById retrieves Errmsg by Id. Returns error if
// Id doesn't exist
func (this *Errmsg) GetErrmsgByCode(code int) (msg string) {
	ormer := orm.NewOrm()

	if code == 0 {
		msg = "OK"

		return
	} else {
		msg = "请稍后再试"
	}

	errmsgModel := Errmsg{}
	oneErr := ormer.QueryTable(new(Errmsg)).Filter("code", code).One(&errmsgModel, "msg")
	// 查询失败
	if oneErr != nil {
		return
	}

	// 数据不存在
	if oneErr == orm.ErrNoRows {
		return
	}

	msg = errmsgModel.Msg

	return
}
