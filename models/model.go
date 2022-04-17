package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" //需要包含驱动package
)

// 表的设计
type User struct {
	Id int
	Name string
	Pwd string
}

func init() {
	// 1.设置数据库基本信息
	mysqladmin := beego.AppConfig.String("mysql_user")
	mysqlpwd := beego.AppConfig.String("mysql_password")
	mysqlhost := beego.AppConfig.String("mysql_host")
	mysqldb := beego.AppConfig.String("mysql_dbname")
	//链接数据库
	orm.RegisterDataBase("default", "mysql", mysqladmin+":"+mysqlpwd+"@tcp("+mysqlhost+")/"+mysqldb+"?charset=utf8")
	//2.映射model数据
	orm.RegisterModel(new(User))
	//3.生成表
	orm.RunSyncdb("default", false, true)//db名称，是否强制更新 是否可见
}