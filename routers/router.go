package routers

import (
	"beego_Study/controllers"
	"github.com/astaxie/beego"
)

func init() {
	// 路由: url地址匹配
	// 不设置，默认访问请求对应的函数
    beego.Router("/", &controllers.MainController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/register", &controllers.MainController{}, "get:ShowRegister;post:HandleRegister")
	//注意：实现自定义的get请求方法，将不会访问默认方法(所有都不会访问)
	//beego.Router("/login", &controllers.MainController{}, "get:ShowLogin;post:HandleLogin")

	// 文章相关
	beego.Router("/Article/ShowArticle", &controllers.ArticleController{}, "get:ShowArticleList")
	beego.Router("/Article/AddArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")
	beego.Router("/Article/content", &controllers.ArticleController{}, "get:ShowContent")
	beego.Router("/Article/UpdateArticle", &controllers.ArticleController{}, "get:ShowUpdate;post:HandleUpdate")
	beego.Router("/Article/DeleteArticle", &controllers.ArticleController{}, "get:HandleDelete")

}

