package routers

import (
	"beego_Study/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// 路由: url地址匹配
    beego.Router("/", &controllers.MainController{})
	beego.Router("/abc", &controllers.MainController{})
	beego.Router("/register", &controllers.MainController{}, "get:ShowRegister;post:HandleRegister")
	//注意：实现自定义的get请求方法，将不会访问默认方法(所有都不会访问)
	beego.Router("/login", &controllers.MainController{}, "get:ShowLogin;post:HandleLogin")

}
