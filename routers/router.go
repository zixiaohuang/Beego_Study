package routers

import (
	"beego_Study/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// 路由: url地址匹配
    beego.Router("/", &controllers.MainController{})
	beego.Router("/abc", &controllers.MainController{})
}
