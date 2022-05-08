package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql" //需要包含驱动package
	"time"
)

// 表的设计
type User struct {
	Id int
	Name string `orm:"unique"` // 全局唯一
	Passwd string
	Articles []*Article `orm:"rel(m2m)"` // ManyToMany relation 多对多关系
}

// 文章结构体
type Article struct {
	Id int `orm:"pk;auto"` // pk主键 auto自增
	Title string `orm:"size(20)"`//文章标题
	Content string `orm:"size(500)"`// 文章内容
	Img string `orm:"size(50)"`//图片路径
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"` //发布时间 auto_now_add 第一次保存时才设置时间
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`   //修改时间 auto_now 每次model保存时都会对时间自动更新
	Count int `orm:"null,default(0)"`// 阅读量 默认值为0,允许为空

	//beego没有使用mysql的原生时间戳，而是自行打时间戳
	// mysql时间设置有两种: date datetime beego为time.time
	ArticleType *ArticleType `orm:"rel(fk)"` // RealForeignKey relation
	Users       []*User      `orm:"reverse(many)"` // 跟rel(m2m)对应， mysql会另外创建一个关系表
}

type ArticleType struct {
	Id       int        `orm:"pk;auto"`
	TypeName string     `orm:"size(20)"`
	Articles []*Article `orm:"reverse(many)"` // 跟rel(fk)相对应，形成一对多关系
}

// orm作用：1.能够通过对象操作相应的数据库表 2.能够通过结构体对象生成对应的数据库表

func init() {
	// 1.设置数据库基本信息
	mysqladmin := beego.AppConfig.String("mysql_user")
	mysqlpwd := beego.AppConfig.String("mysql_password")
	mysqlhost := beego.AppConfig.String("mysql_host")
	mysqldb := beego.AppConfig.String("mysql_dbname")
	//链接数据库
	orm.RegisterDataBase("default", "mysql", mysqladmin+":"+mysqlpwd+"@tcp("+mysqlhost+")/"+mysqldb+"?charset=utf8")
	//2.映射model数据
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	//3.生成表
	orm.RunSyncdb("default", false, true)//db名称，是否强制更新 是否可见
}