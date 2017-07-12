package main

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

type Student struct {
	Id    int // primary key
	Name  string
	Age   int
	Sex   string
	Score float32
	Addr  string
}

func init() {
	orm.RegisterModel(new(Student))

	orm.RegisterDriver("sqlite", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "./test.db")
	orm.RunSyncdb("default", false, true)
}

func main() {
	o := orm.NewOrm()
	o.Using("default") // 默认使用 default，你可以指定为其他数据库

	stu := new(Student)
	stu.Age = 30
	stu.Name = "allen"
	stu.Sex = "male"
	stu.Score = 88
	stu.Addr = "Beijing, China"

	fmt.Println(o.Insert(stu))
	var r orm.RawSeter
	r = o.Raw("select * from student;")
	var s []Student
	r.QueryRows(&s)
	fmt.Println(s)

	var qs orm.QuerySeter
	qs = o.QueryTable("student")
	fmt.Println(qs.Offset(1).All(&s))
	fmt.Println(s)
}
