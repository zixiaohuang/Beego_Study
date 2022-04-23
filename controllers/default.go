package controllers

import (
	"beego_Study/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type MainController struct {
	beego.Controller // 继承
}

//// 对应get请求
//func (c *MainController) Get() {
//	// Website Email传给tpl模版
//	c.Data["Website"] = "beego.me"
//	c.Data["Email"] = "astaxie@gmail.com"
//	c.TplName = "index.tpl" //模版
//}

func (c *MainController) Get() {
	c.TplName = "login.html" //模版
}

func (c *MainController) ShowRegister() {
	// 插入
	//// 1.有orm对象
	//o := orm.NewOrm()
	//// 2.有一个要插入数据的结构体对象
	//user := models.User{}
	//// 3.对结构体对象赋值
	//user.Name = "zixiao"
	//user.Pwd = "1234"
	//// 4.插入
	//_, err := o.Insert(&user)
	//if err != nil {
	//	logs.Error("插入失败", err)
	//	return
	//}

	//// 查询
	//// 1.有orm对象
	//o := orm.NewOrm()
	//// 2.查询对象
	//user := models.User{}
	//// 3.确定查询对象字段值
	////user.Id = 1
	//// 4.查询
	////err := o.Read(&user)
	//
	//user.Name = "zixiao"
	//err := o.Read(&user, "Name")
	//if err != nil {
	//	logs.Error("查询失败", err)
	//	return
	//}
	//logs.Info("查询成功", user)

	//// 更新
	//// 1.orm对象
	//o := orm.NewOrm()
	//// 2.需要更新的结构体对象
	//user := models.User{}
	//// 3.查询到需要更新的数据
	//user.Id = 1
	//err := o.Read(&user)
	//if err != nil {
	//	logs.Error("查询失败", err)
	//	return
	//}
	//
	//// 4.给数据重新赋值
	//user.Pwd = "234"
	//// 5.更新
	//_, err = o.Update(&user)
	//if err != nil {
	//	logs.Error("更新失败", err)
	//	return
	//}

	// 删除
	//// 1. ORM对象
	//o := orm.NewOrm()
	//// 2. 删除对象
	//user := models.User{}
	//// 3. 指定删除那一条数据
	//user.Id = 1
	//// 4.删除
	//_, err := o.Delete(&user)
	//if err != nil {
	//	logs.Error("删除失败", err)
	//	return
	//}

	c.TplName = "register.html"
}

func (c *MainController)HandleRegister(){
	// 1.注册业务 拿到数据
	userName := c.GetString("userName")
	pwd := c.GetString("password")
	logs.Info(userName, pwd)

	// 2.对数据进行校验
	if userName == "" || pwd == "" {
		logs.Error("数据不能为空")
		c.Redirect("/register", 302) //第一个参数：url地址 第二个参数： http状态码，302重定向
		return
	}
	// 3.插入数据库
	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	user.Passwd = pwd
	_, err := o.Insert(&user)
	if err != nil {
		logs.Info("插入数据库失败")
		c.Redirect("/register", 302)
		return
	}
	// 4.返回登陆界面
	// c.TplName 和 c.Redirect的区别:
	// c.TplName 指定视图文件，同时可以给这个视图传递一些数据
	// c.Redirect 跳转，不能传递数据 数据快
	// 1xx:请求已经被接受，需要继续发送请求 2xx:请求成功 3xx:请求资源被转移，请求被转接
	// 4xx:请求失败，客户端 5xx: 服务器错误
	c.Redirect("/login", 302)
}

func (c *MainController)ShowLogin() {
	c.TplName = "login.html"
}

func (c *MainController)HandleLogin(){
	//c.Ctx.WriteString("这是登陆的Post请求")
	// 登陆: 1.拿到数据
	userName := c.GetString("userName")
	pwd := c.GetString("password")
	// 2.判断数据是否合法
	if userName == ""|| pwd == "" {
		logs.Info("输入数据不合法")
		c.TplName="login.html"
		return
	}

	// 3.查询账号密码是否正确
	o := orm.NewOrm()
	user := models.User{}
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		logs.Info("查询失败")
		c.TplName = "login.html"
		return
	}

	// 密码校验
	if user.Passwd != pwd {
		logs.Info("密码不正确")
		c.TplName = "login.html"
		return
	}

	// 4.跳转
	c.Redirect("/Article/ShowArticle", 302)
	//c.Ctx.WriteString("欢迎您，登陆成功")
}

