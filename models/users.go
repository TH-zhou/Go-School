package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type Users struct {
	Id               int    `orm:"column(id);auto"`
	Corpid           int    `orm:"column(corpid)" description:"商户id"`
	Userid           int    `orm:"column(userid)" description:"用户id"`
	Name             string `orm:"column(name);size(11);null" description:"名称"`
	Tel              string `orm:"column(tel);size(11)" description:"手机号"`
	Photo            string `orm:"column(photo);size(100);null" description:"用户头像"`
	Did              string `orm:"column(did);size(20);null" description:"部门id 多个部门用逗号隔开"`
	DidParent        string `orm:"column(did_parent);size(100);null" description:"用户所在部门下所有子部门，逗号分隔"`
	PostId           string `orm:"column(post_id);size(20);null" description:"岗位id 多个岗位用户逗号隔开"`
	RoleId           int    `orm:"column(role_id);null" description:"角色id"`
	IsAdmin          int8   `orm:"column(is_admin);null" description:"1为超级管理员 0不是"`
	Score            int    `orm:"column(score);null" description:"积分"`
	Duration         int    `orm:"column(duration);null" description:"在线时长"`
	ExamPassingRate  string `orm:"column(exam_passing_rate);size(4);null" description:"考试通过率 0表示0% 50表示50%"`
	Status           int8   `orm:"column(status);null" description:"状态 1：正常 0：暂停 2：删除"`
	Entrytime        int    `orm:"column(entrytime);null" description:"入职时间"`
	Createtime       int    `orm:"column(createtime);null" description:"入职时间"`
	ResetPwdCount    int8   `orm:"column(reset_pwd_count);null" description:"每天重置密码次数"`
	LastResetPwdTime int    `orm:"column(last_reset_pwd_time);null" description:"上次重置密码时间"`
}

func (t *Users) TableName() string {
	return "users"
}

func init() {
	orm.RegisterModel(new(Users))
}

// AddUsers insert a new Users into database and returns
// last inserted Id on success.
func AddUsers(m *Users) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUsersById retrieves Users by Id. Returns error if
// Id doesn't exist
func GetUsersById(id int) (v *Users, err error) {
	o := orm.NewOrm()
	v = &Users{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUsers retrieves all Users matches certain condition. Returns empty list if
// no records exist
func GetAllUsers(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Users))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Users
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateUsers updates Users by Id and returns error if
// the record to be updated doesn't exist
func UpdateUsersById(m *Users) (err error) {
	o := orm.NewOrm()
	v := Users{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUsers deletes Users by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUsers(id int) (err error) {
	o := orm.NewOrm()
	v := Users{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Users{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
