package controllers

import (
	"github.com/astaxie/beego"
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
	c.Data["data"] = "今天吃饺子"
	c.TplName = "test.html"
}

func (c *MainController)Post(){
	c.Data["data"] = "今天吃面条"
	c.TplName = "test.html"
}
