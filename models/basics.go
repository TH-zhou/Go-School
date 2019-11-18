package models

import (
	"github.com/astaxie/beego/orm"
)

type Basics struct {
	Id         int    `orm:"column(id);auto"`
	Corpid     int    `orm:"column(corpid)" description:"商户id"`
	Type       string `orm:"column(type);size(20)" description:"类型"`
	Alias      string `orm:"column(alias);size(20);null" description:"名称"`
	Value      string `orm:"column(value);size(20);null" description:"类型值"`
	Sort       int    `orm:"column(sort);null" description:"排序"`
	IsDelete   int8   `orm:"column(is_delete);null" description:"状态 0未删除 1已删除"`
	Createtime int    `orm:"column(createtime);null" description:"创建时间"`
}

func (t *Basics) TableName() string {
	return "basics"
}

func init() {
	orm.RegisterModel(new(Basics))
}

// 获取该类型数据
func (this *Basics) GetBasicsTypeByvalue(corpid int, basesType string, value int) (*Basics, error) {
	ormer := orm.NewOrm()
	cond := orm.NewCondition()

	var basic Basics
	if oneErr := ormer.QueryTable(new(Basics)).Filter("type", basesType).Filter("value", value).SetCond(cond.AndCond(cond.And("corpid", 0).Or("corpid", corpid))).One(&basic); oneErr != nil {
		return nil, oneErr
	}

	return &basic, nil
}
